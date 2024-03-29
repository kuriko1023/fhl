package main

import (
	"bufio"
	"container/heap"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
)

// 一篇诗词
type Article struct {
	Id      int      // 编号
	Title   string   // 标题
	Dynasty string   // 朝代
	Author  string   // 作者
	Content []string // 内容，由标点拆分出的句子列表
}

// 所有诗词的列表
// 每一项的 Id 等于此列表中的下标
// 后续建立了一个 LRU 缓存，某篇目不在缓存中时，对应项为 nil
var articles []*Article

// 每个篇目在文件中的偏移值
var articleOffset []int64

// 名篇的列表
var hotArticles []*Article

// 所有高频词组成的列表，单字和双字分开，各自按频率降序排序
var hotWords1 []string
var hotWords2 []string

// 高频词的出现频数，单字和双字分开
type RunePair struct{ A, B rune }

var hotWords1Count map[rune]int
var hotWords2Count map[RunePair]int

// 全部都是高频字的句子
var allHotSentences [][]string

const ALL_HOT_LEN_MIN = 5
const ALL_HOT_LEN_MAX = 9

//高频词组合（单字＆单字／单字＆双字／双字＆双字）在诗句中出现的频次
//用于优化谜之飞花令的题目选择
var hotWordsFreq map[string]int

// 返回一句诗词中的所有高频词，按出现次数降序排序；若无，返回空列表
// 单字高频词与双字高频词分别返回，第一个返回值为单字
func getHotWords(sentence string) ([]string, []string) {
	count1 := []KVPair{}
	count2 := []KVPair{}

	// 单字词
	s := []rune(sentence)
	for i, c := range s {
		if n, has := hotWords1Count[c]; has {
			count1 = append(count1, KVPair{string(c), n})
		}
		if i < len(s)-1 {
			if n, has := hotWords2Count[RunePair{c, s[i+1]}]; has {
				count2 = append(count2, KVPair{
					string(c) + string(s[i+1]),
					n,
				})
			}
		}
	}

	sort.Sort(byValueDesc(count1))
	sort.Sort(byValueDesc(count2))

	singleHotWords := []string{}
	doubleHotWords := []string{}
	for _, p := range count1 {
		singleHotWords = append(singleHotWords, p.string)
	}
	for _, p := range count2 {
		doubleHotWords = append(doubleHotWords, p.string)
	}
	return singleHotWords, doubleHotWords
}

func randomSample(n, count int) []int {
	picked := map[int]struct{}{}
	ret := []int{}
	for i := n - count; i < n; i++ {
		x := rand.Intn(i + 1)
		if _, dup := picked[x]; dup {
			x = i
		}
		picked[x] = struct{}{}
		ret = append(ret, x)
	}
	return ret
}

// 生成普通飞花题目备选列表
// 返回 count1 个不重复的单字与 count2 个不重复的双字
func generateA(count1, count2 int) []string {
	ret := []string{}

	// 单字
	if count1 > 0 {
		n := len(hotWords1) / 2
		for _, i := range randomSample(n, count1-1) {
			ret = append(ret, hotWords1[i])
		}
		ret = append(ret, hotWords1[rand.Intn(n)+n])
	}

	// 双字
	if count2 > 0 {
		n := len(hotWords2) / 2
		for _, i := range randomSample(n, count2-1) {
			ret = append(ret, hotWords2[i])
		}
		ret = append(ret, hotWords2[rand.Intn(n)+n])
	}

	return ret
}

// 生成多字飞花题目，返回一句长度为 length 的句子
// 需确保 ALL_HOT_LEN_MIN <= length <= ALL_HOT_LEN_MAX
func generateB(length int) string {
	collection := allHotSentences[length-ALL_HOT_LEN_MIN]
	return collection[rand.Intn(len(collection))]
}

// 生成超级飞花题目，返回长度为 sizeLeft 和 sizeRight 的字符串列表
func generateC(sizeLeft, sizeRight int) ([]string, []string) {
	var left []string
	var right []string

	leftTopPart := sizeLeft
	leftBottomPart := 0
	if sizeLeft >= 3 {
		leftTopPart -= 1
		leftBottomPart += 1
	}
	n := len(hotWords1)

	// 左侧的一个字最多对应右侧的多少字
	singleWordLimit := sizeRight/sizeLeft + 2
	word2Limit := sizeRight / 5

	for len(right) < sizeRight {
		// 所有已选字的集合
		all := map[string]struct{}{}

		// 先选取 n 个高频字
		left = []string{}
		for _, i := range randomSample(n/4, leftTopPart) {
			left = append(left, hotWords1[i])
		}
		for _, i := range randomSample(n-n/4, leftBottomPart) {
			left = append(left, hotWords1[i+n/4])
		}
		for _, c := range left {
			all[c] = struct{}{}
		}

		// 寻找包含任意选出字的名篇句子，并从中找出一些高频词
		// 如果性能不足可以后续优化
		right = []string{}
		// 记录 left 中每个字已经对应了 right 中多少个字
		count := make([]int, len(left))
		// 有多少个双字
		word2Count := 0
		// 先打乱
		perm := rand.Perm(len(hotArticles))
		for _, id := range perm {
			article := hotArticles[id]
			sentenceIndex := -1
			for j, t := range left {
				if count[j] >= singleWordLimit {
					continue
				}
				for i, s := range article.Content {
					if len([]rune(s)) >= 4 && strings.Contains(s, t) {
						count[j] += 1
						sentenceIndex = i
						break
					}
				}
				if sentenceIndex != -1 {
					break
				}
			}
			if sentenceIndex == -1 {
				continue
			}
			h1, h2 := getHotWords(article.Content[sentenceIndex])
			// 先考虑双字词，如有且不与单字重复，则加入
			word2Valid := false
			if word2Count < word2Limit {
				for _, word := range h2 {
					runes := []rune(word)
					_, picked0 := all[string(runes[0])]
					_, picked1 := all[string(runes[1])]
					_, picked := all[word]
					if !picked0 && !picked1 && !picked {
						right = append(right, word)
						all[string(runes[0])] = struct{}{}
						all[string(runes[1])] = struct{}{}
						all[word] = struct{}{}
						word2Valid = true
						word2Count++
						break
					}
				}
			}
			if !word2Valid {
				// 随机选取一个常见字
				// 先打乱
				for i := range h1 {
					j := rand.Intn(i + 1)
					h1[i], h1[j] = h1[j], h1[i]
				}
				for _, word := range h1 {
					if _, picked := all[word]; !picked {
						right = append(right, word)
						all[word] = struct{}{}
						break
					}
				}
			}
			if len(right) == sizeRight {
				break
			}
		}
	}

	sort.Slice(right, func(i, j int) bool {
		return len(right[i]) < len(right[j])
	})

	return left, right
}

// 生成 2 个长度为 n 的关键词组作为谜之飞花题目
//其中第二个关键词组全部由双字高频词构成,第一个关键词组由单字或双字高频词构成
func generateD(n int) ([]string, []string) {
	var hotWordsList1 []string
	var hotWordsList2 []string

	count := 0
	length := len(hotArticles)
	//储存已被选择过的诗
	var tmp map[int]int
	tmp = make(map[int]int)

	for count < n {
		i := rand.Intn(length)
		//该诗已被选择过
		if _, ok := tmp[i]; ok {
			continue
		}
		tmp[i] = 1

		content := hotArticles[i].Content
		//随机选择某一句
		j := rand.Intn(len(content))

		sHotWords, dHotWords := getHotWords(content[j])
		//由诗句获得的单字高频词组、双字高频词组都为空
		if len(sHotWords) == 0 && len(dHotWords) == 0 {
			continue
		}

		freqMax := 0
		s1 := 0
		s2 := 0
		flag := 0

		k := rand.Intn(6)
		//有0.6的概率填入（若可行）两个单字词
		if k >= 4 && len(sHotWords) > 1 {
			for i := 0; i < len(sHotWords); i++ {
				for j := i + 1; j < len(sHotWords); j++ {
					if freq, ok := hotWordsFreq[sHotWords[i]+sHotWords[j]]; ok {
						//高频词组合频次大于50方记入; 取最大值
						if freq >= 50 && freq > freqMax {
							s1 = i
							s2 = j
							flag = 1
						}
					}
				}
			}
			if flag == 1 {
				hotWordsList1 = append(hotWordsList1, sHotWords[s1])
				hotWordsList2 = append(hotWordsList2, sHotWords[s2])
			}
		}
		//有0.4的概率填入（若可行）单字词 + 双字词
		if k < 4 && len(sHotWords) > 0 && len(dHotWords) > 0 {
			for i := 0; i < len(sHotWords); i++ {
				for j := 0; j < len(dHotWords); j++ {
					//确保单字词不包含在双字词中
					if !strings.Contains(dHotWords[j], sHotWords[i]) {
						if freq, ok := hotWordsFreq[sHotWords[i]+dHotWords[j]]; ok {
							//高频词组合频次大于50方记入; 取最大值
							if freq >= 50 && freq > freqMax {
								s1 = i
								s2 = j
								flag = 1
							}
						}
					}
				}
			}
			if flag == 1 {
				hotWordsList1 = append(hotWordsList1, sHotWords[s1])
				hotWordsList2 = append(hotWordsList2, dHotWords[s2])
			}
		}

		if flag == 0 {
			continue
		}

		//计数加1
		count++
		// fmt.Println(content[j])
	}

	sort.Strings(hotWordsList1)
	sort.Strings(hotWordsList2)

	return hotWordsList1, hotWordsList2
}

// 排序用比较器
type KVPair struct {
	string
	int
}
type byValueDesc []KVPair

func (s byValueDesc) Len() int {
	return len(s)
}
func (s byValueDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byValueDesc) Less(i, j int) bool {
	return s[i].int > s[j].int
}

func furtherInit() {
	hotWordsFreq = make(map[string]int)

	// 初始化高频词组合频次表，令其包括所有单字&单字、单字&双字、双字&双字的组合
	for i := 0; i < len(hotWords1); i++ {
		for j := i + 1; j < len(hotWords1); j++ {
			hotWordsFreq[hotWords1[i]+hotWords1[j]] = 0
		}
		for j := 0; j < len(hotWords2); j++ {
			//若单字词为双字词的一部分,忽略
			if !strings.Contains(hotWords2[j], hotWords1[i]) {
				hotWordsFreq[hotWords1[i]+hotWords2[j]] = 0
			}
		}
	}

	for _, article := range articles {
		content := article.Content
		for i := 0; i < len(content)-1; i++ {
			//拼接两句为一“联”
			sentence := content[i] + content[i+1]
			sHotWords, dHotWords := getHotWords(sentence)

			//根据每一“联”诗词中的高频词，更新高频词组合频次表
			for k := 0; k < len(sHotWords); k++ {
				for j := k + 1; j < len(sHotWords); j++ {
					hotWordsFreq[sHotWords[k]+sHotWords[j]]++
				}
				for j := 0; j < len(dHotWords); j++ {
					if !strings.Contains(dHotWords[j], sHotWords[k]) {
						hotWordsFreq[sHotWords[k]+dHotWords[j]]++
					}
				}
			}
		}
	}
	// for k, v := range hotWordsFreq {
	// 	fmt.Println(k,v)
	// }
	// fmt.Println("finish")

	for k, v := range hotWordsFreq {
		if v < 50 {
			delete(hotWordsFreq, k)
		}
	}
}

var datasetFile *os.File

var gobValues = []interface{}{
	&articleOffset,
	&hotWords1,
	&hotWords2,
	&hotWords1Count,
	&hotWords2Count,
	&allHotSentences,
	&hotWordsFreq,
	&hotArticles, // TODO: 这好吗？这不好
}

var precalFile *os.File
var errCorrOffset int64
var errCorrNumRecords int64

func loadPrecal() error {
	file, err := os.Open("../dataset/2c-precal.bin")
	if err != nil {
		return err
	}

	// 解码 Gob
	decoder := gob.NewDecoder(file)
	for _, v := range gobValues {
		if err := decoder.Decode(v); err != nil {
			file.Close()
			return err
		}
	}

	// 创建空的 articles 数组
	articles = make([]*Article, len(articleOffset))

	// 纠错数据
	offs, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		file.Close()
		return err
	}
	errCorrOffset = offs

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	errCorrNumRecords = (stat.Size() - offs) / RECORD_W

	precalFile = file
	return nil
}

func savePrecalGob() error {
	file, err := os.Create("../dataset/2c-precal.bin")
	if err != nil {
		return err
	}

	// 保存 Gob
	encoder := gob.NewEncoder(file)
	for _, v := range gobValues {
		if err := encoder.Encode(v); err != nil {
			file.Close()
			return err
		}
	}

	// 纠错数据将在之后保存
	precalFile = file
	return nil
}

func savePrecalErrCorr(x []ErrCorrRecord) error {
	w := bufio.NewWriter(precalFile)

	count := 0
	for i, rec := range x {
		if i == 0 || rec != x[i-1] {
			count++
			if err := writeErrCorrRecord(w, rec); err != nil {
				return err
			}
		}
	}

	w.Flush()
	println(count)
	return nil
}

func parseArticle(id int, s string) (*Article, string) {
	fields := strings.SplitN(s, "\t", 5)
	if len(fields) < 5 {
		panic("Incorrect dataset format")
	}

	return &Article{
		Id:      id,
		Title:   fields[3],
		Dynasty: fields[1],
		Author:  fields[2],
		Content: strings.Split(fields[4], "/"),
	}, fields[0]
}

// 读入数据集，填充所有全局变量
func initDataset() {
	rand.Seed(1023)

	file, err := os.Open("../dataset/2b-dedup.txt")
	if err != nil {
		panic(err)
	}
	datasetFile = file

	if loadPrecal() == nil {
		initArticleCache()
		return
	}

	articles = []*Article{}
	articleOffset = []int64{}
	hotArticles = []*Article{}
	hotWords1 = []string{}
	hotWords2 = []string{}

	hotWords1Count = map[rune]int{}
	hotWords2Count = map[RunePair]int{}

	i := 0
	sc := bufio.NewScanner(file)
	offs := int64(0)
	p := 0
	q := 0
	t := 0
	for sc.Scan() {
		prevOffs := offs
		offs += int64(len(sc.Text())) + 1

		// 随机抽取十分之一
		i++
		/*if sc.Text()[0] != '!' && i%10 != 0 {
			continue
		}*/

		// 将篇目加入列表
		article, flag := parseArticle(len(articles), sc.Text())
		articleOffset = append(articleOffset, prevOffs)
		articles = append(articles, article)
		weight := 1
		if flag == "!" {
			// 名篇
			hotArticles = append(hotArticles, article)
			weight = 10
		}

		for _, s := range article.Content {
			n := len([]rune(s))
			p += 1
			q += n
			t += n*(n+1)/2 + 1
		}

		// 若不是重复篇目，则计入高频词
		if flag != " " {
			for _, s := range article.Content {
				s := []rune(s)
				for i, c := range s {
					hotWords1Count[c] += weight
					if i < len(s)-1 {
						hotWords2Count[RunePair{c, s[i+1]}] += weight
					}
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("dataset: %d articles\n", len(articles))
	fmt.Printf("%d, %d, %d\n", p, q, t)

	hotWords1List := byValueDesc{}
	hotWords2List := byValueDesc{}
	for k, v := range hotWords1Count {
		hotWords1List = append(hotWords1List, KVPair{string(k), v})
	}
	for k, v := range hotWords2Count {
		hotWords2List = append(hotWords2List, KVPair{
			string(k.A) + string(k.B),
			v,
		})
	}
	sort.Sort(hotWords1List)
	sort.Sort(hotWords2List)

	fmt.Println("高频单字，每行五十个")
	for i := 0; i < 400; i++ {
		fmt.Print(hotWords1List[i].string)
		if (i+1)%50 == 0 {
			fmt.Println()
		}
		hotWords1 = append(hotWords1, hotWords1List[i].string)
	}
	fmt.Println("高频双字，每行二十个")
	for i := 0; i < 200; i++ {
		fmt.Print(hotWords2List[i].string, " ")
		if (i+1)%20 == 0 {
			fmt.Println()
		}
		hotWords2 = append(hotWords2, hotWords2List[i].string)
	}

	// 删去非高频词
	for i := len(hotWords1); i < len(hotWords1List); i++ {
		runes := []rune(hotWords1List[i].string)
		delete(hotWords1Count, runes[0])
	}
	for i := len(hotWords2); i < len(hotWords2List); i++ {
		runes := []rune(hotWords2List[i].string)
		delete(hotWords2Count, RunePair{runes[0], runes[1]})
	}

	// 找出仅由不重复的高频字组成的句子
	allHotSentences = make([][]string, ALL_HOT_LEN_MAX-ALL_HOT_LEN_MIN+1)
	allHotSentencesSet := make(map[string]struct{})
	for i := range allHotSentences {
		allHotSentences[i] = []string{}
	}
	for _, article := range articles {
		for _, s := range article.Content {
			runes := []rune(s)
			if len(runes) < ALL_HOT_LEN_MIN || len(runes) > ALL_HOT_LEN_MAX {
				continue
			}
			// 确认每个字是否是高频字
			minCount := hotWords1List[len(hotWords1)/2-1].int
			if len(runes) >= ALL_HOT_LEN_MAX-2 {
				minCount = hotWords1List[len(hotWords1)-1].int
			}
			allHot := true
			for _, r := range runes {
				if n, has := hotWords1Count[r]; !has || n < minCount {
					allHot = false
					break
				}
			}
			if allHot {
				i := len(runes) - ALL_HOT_LEN_MIN
				// 确认是否有重复字
				// 因为 runes 至多只有九个元素，所以用朴素算法
				for i, r := range runes {
					for j := 0; j < i; j++ {
						if runes[j] == r {
							goto out
						}
					}
				}
				// 去重
				if _, existing := allHotSentencesSet[s]; existing {
					goto out
				}
				// 通过检查，加入列表
				allHotSentences[i] = append(allHotSentences[i], s)
				allHotSentencesSet[s] = struct{}{}
			out:
			}
		}
	}
	fmt.Println("仅由不重复的高频字组成的句子")
	for i, c := range allHotSentences {
		fmt.Printf("%d 字：%d 句\n", i+ALL_HOT_LEN_MIN, len(c))
	}

	furtherInit()
	if err := savePrecalGob(); err != nil {
		panic(err)
	}

	initErrCorr()

	precalFile.Close()
	if err := loadPrecal(); err != nil {
		panic(err)
	}

	for i, _ := range articles {
		articles[i] = nil
	}
	initArticleCache()
}

// 篇目的 LRU cache
type CacheEntry struct {
	Id int
	Ts uint64
}
type CacheEntryHeap []CacheEntry

func (h CacheEntryHeap) Len() int            { return len(h) }
func (h CacheEntryHeap) Less(i, j int) bool  { return h[i].Ts < h[j].Ts }
func (h CacheEntryHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *CacheEntryHeap) Push(x interface{}) { *h = append(*h, x.(CacheEntry)) }
func (h *CacheEntryHeap) Pop() interface{} {
	x := (*h)[len(*h)-1]
	*h = (*h)[0 : len(*h)-1]
	return x
}

// 元素的 a: 篇目编号；b: 上次访问的时间戳
var articleCache = &CacheEntryHeap{}
var articleCacheTimestamp = uint64(0)
var articleActualTimestamp []uint64

var datasetFileReader = bufio.NewReaderSize(nil, 512)

func initArticleCache() {
	heap.Init(articleCache)
	articleActualTimestamp = make([]uint64, len(articles))
}

func getArticle(id int) *Article {
	articleCacheTimestamp++ // 应该不会溢出吧
	articleActualTimestamp[id] = articleCacheTimestamp
	// println("accessing", id)

	if a := articles[id]; a != nil {
		return a
	}

	// 读入篇目
	datasetFile.Seek(articleOffset[id], os.SEEK_SET)
	datasetFileReader.Reset(datasetFile)
	s, err := datasetFileReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s = s[:len(s)-1] // 去掉换行符
	article, _ := parseArticle(id, s)

	// 加入缓存
	// 若缓存已满，则移除最久未使用的一个
	if len(*articleCache) >= Config.ArticleCache {
		for {
			elm := heap.Pop(articleCache).(CacheEntry)
			ts := articleActualTimestamp[elm.Id]
			if ts != elm.Ts {
				// 向堆中加入正确的值
				heap.Push(articleCache, CacheEntry{elm.Id, ts})
				continue
			}
			// println("removing", elm.Id)
			articles[elm.Id] = nil
			break
		}
	}

	heap.Push(articleCache, CacheEntry{id, articleCacheTimestamp})
	articles[id] = article

	return article
}

// 对称删除纠错算法

type HashType uint64
type ArticleIdxType uint32
type ContentIdxType uint32
type ErrCorrRecord struct {
	Hash       HashType
	ArticleIdx ArticleIdxType
	ContentIdx ContentIdxType
}

const HASH_W = 5
const ART_IDX_W = 3
const CON_IDX_W = 2
const RECORD_W = HASH_W + ART_IDX_W + CON_IDX_W

func hash(rs []rune) HashType {
	h := HashType(0)
	for _, r := range rs {
		if r > 0 {
			h = h*100003 + HashType(r)
		}
	}
	return h % (1 << (HASH_W * 8))
}

func readErrCorrRecord(index int64) ErrCorrRecord {
	buf := [RECORD_W]byte{}
	precalFile.ReadAt(buf[:], errCorrOffset+index*RECORD_W)
	rec := ErrCorrRecord{0, 0, 0}
	for i, b := range buf[0:HASH_W] {
		rec.Hash += (HashType(b) << (i * 8))
	}
	for i, b := range buf[HASH_W : HASH_W+ART_IDX_W] {
		rec.ArticleIdx += (ArticleIdxType(b) << (i * 8))
	}
	for i, b := range buf[HASH_W+ART_IDX_W:] {
		rec.ContentIdx += (ContentIdxType(b) << (i * 8))
	}
	return rec
}

func writeErrCorrRecord(w *bufio.Writer, rec ErrCorrRecord) error {
	buf := [RECORD_W]byte{}
	for i := 0; i < HASH_W; i++ {
		buf[i] = byte(rec.Hash >> (i * 8))
	}
	for i := 0; i < ART_IDX_W; i++ {
		buf[HASH_W+i] = byte(rec.ArticleIdx >> (i * 8))
	}
	for i := 0; i < CON_IDX_W; i++ {
		buf[HASH_W+ART_IDX_W+i] = byte(rec.ContentIdx >> (i * 8))
	}
	_, err := w.Write(buf[:])
	return err
}

func forEachPossibleErrHash(s string, fn func(h HashType) bool) {
	rs := []rune(s)
	if fn(hash(rs)) {
		return
	}
	for i, r := range rs {
		rs[i] = -1
		if fn(hash(rs)) {
			return
		}
		for j, r := range rs[:i] {
			rs[j] = -1
			if fn(hash(rs)) {
				return
			}
			rs[j] = r
		}
		rs[i] = r
	}
}

func initErrCorr() {
	x := []ErrCorrRecord{}
	for i, article := range articles {
		for j, s := range article.Content {
			forEachPossibleErrHash(s, func(h HashType) bool {
				x = append(x, ErrCorrRecord{
					Hash:       h,
					ArticleIdx: ArticleIdxType(i),
					ContentIdx: ContentIdxType(j),
				})
				return false
			})
		}
	}
	println(len(x))
	sort.Slice(x, func(i, j int) bool {
		return x[i].Hash < x[j].Hash
	})

	if err := savePrecalErrCorr(x); err != nil {
		panic(err)
	}
}

// 在纠错数据库中查找某个 hash 值
// 返回 >= 此 hash 的最小记录位置，即 lower_bound
func lookupErrCorr(x HashType) int64 {
	lo := int64(-1)
	hi := errCorrNumRecords
	for lo < hi-1 {
		mid := (lo + hi) / 2
		rec := readErrCorrRecord(mid)
		if rec.Hash < x {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi
}

// 检查句子是否在诗词库中
// 返回：(是否完全一致, 篇目编号, 句子下标)
// 找不到时，返回的篇目编号与句子下标均为 -1
func lookupText(text []string) (bool, int, int) {
	// 找到第一个至少含四字的子句；若无则选择第一个
	pivot := 0
	for i, s := range text {
		if len([]rune(s)) >= 4 {
			pivot = i
			break
		}
	}

	// 是否有接近
	bestDist := 3 // 最大允许的距离 + 1
	bestArticle := -1
	bestContent := -1

	forEachPossibleErrHash(text[pivot], func(h HashType) bool {
		index := lookupErrCorr(h)
		for index < errCorrNumRecords {
			rec := readErrCorrRecord(index)
			if rec.Hash != h {
				break
			}

			article := getArticle(int(rec.ArticleIdx))
			i := int(rec.ContentIdx) - pivot
			if i >= 0 && i+len(text) <= len(article.Content) {
				// 检查两段文字是否相同或接近
				templ := article.Content[i : i+len(text)]
				totalDist := 0
				for j, s := range text {
					dist := levenshtein.ComputeDistance(s, templ[j])
					totalDist += dist
					if totalDist >= bestDist {
						break
					}
				}
				if totalDist < bestDist {
					bestDist = totalDist
					bestArticle = int(rec.ArticleIdx)
					bestContent = i
					if totalDist == 0 {
						return true
					}
				}
			}

			index++
		}
		return false
	})

	return (bestDist == 0), bestArticle, bestContent
}
