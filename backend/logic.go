package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
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
var articles []Article

// 名篇的列表
var hotArticles []Article

// 所有高频词组成的列表，单字和双字分开，各自按频率降序排序
var hotWords1 []string
var hotWords2 []string

// 高频词的出现频数，单字和双字分开
type RunePair struct{ a, b rune }
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
					string(c) + string(s[i + 1]),
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
	n := len(hotWords1) / 2
	for _, i := range randomSample(n, count1-1) {
		ret = append(ret, hotWords1[i])
	}
	ret = append(ret, hotWords1[rand.Intn(n)+n])

	// 双字
	n = len(hotWords2) / 2
	for _, i := range randomSample(n, count2-1) {
		ret = append(ret, hotWords2[i])
	}
	ret = append(ret, hotWords2[rand.Intn(n)+n])

	return ret
}

// 生成多字飞花题目，返回一句长度为 length 的句子
// 需确保 ALL_HOT_LEN_MIN <= length <= ALL_HOT_LEN_MAX
func generateB(length int) string {
	collection := allHotSentences[length - ALL_HOT_LEN_MIN]
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
					if freq, ok := hotWordsFreq[sHotWords[i] + sHotWords[j]]; ok{
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
						if freq, ok := hotWordsFreq[sHotWords[i] + dHotWords[j]]; ok{
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
		fmt.Println(content[j])
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
			hotWordsFreq[hotWords1[i] + hotWords1[j]] = 0
		}
		for j := 0; j < len(hotWords2); j++ {
			//若单字词为双字词的一部分,忽略
			if !strings.Contains(hotWords2[j], hotWords1[i]) {
				hotWordsFreq[hotWords1[i] + hotWords2[j]] = 0
			}
		}
	}

	for _, article := range articles {
		content := article.Content
		for i := 0; i < len(content) - 1; i++ {
			//拼接两句为一“联”
			sentence := content[i] + content[i + 1]
			sHotWords, dHotWords := getHotWords(sentence)

			//根据每一“联”诗词中的高频词，更新高频词组合频次表
			for k := 0; k < len(sHotWords); k++ {
				for j := k + 1; j < len(sHotWords); j++ {
					hotWordsFreq[sHotWords[k] + sHotWords[j]]++
				}
				for j := 0; j < len(dHotWords); j++ {
					if !strings.Contains(dHotWords[j], sHotWords[k]) {
						hotWordsFreq[sHotWords[k] + dHotWords[j]]++
					}
				}
			}
		}
	}
	// for k, v := range hotWordsFreq {
	// 	fmt.Println(k,v)
	// }
	// fmt.Println("finish")
}

// 读入数据集，填充所有全局变量
func initDataset() {
	articles = []Article{}
	hotArticles = []Article{}
	hotWords1 = []string{}
	hotWords2 = []string{}

	file, err := os.Open("dataset.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hotWords1Count = map[rune]int{}
	hotWords2Count = map[RunePair]int{}

	i := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		// 随机抽取十分之一
		i++
		if sc.Text()[0] != '!' && i%10 != 0 {
			continue
		}

		fields := strings.SplitN(sc.Text(), "\t", 5)
		if len(fields) < 5 {
			panic("Incorrect dataset format")
		}

		// 将篇目加入列表
		article := Article{
			Id:      len(articles),
			Title:   fields[3],
			Dynasty: fields[1],
			Author:  fields[2],
			Content: strings.Split(fields[4], "/"),
		}
		articles = append(articles, article)
		weight := 1
		if fields[0] == "!" {
			// 名篇
			hotArticles = append(hotArticles, article)
			weight = 10
		}

		// 若不是重复篇目，则计入高频词
		if fields[0] != " " {
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

	hotWords1List := byValueDesc{}
	hotWords2List := byValueDesc{}
	for k, v := range hotWords1Count {
		hotWords1List = append(hotWords1List, KVPair{string(k), v})
	}
	for k, v := range hotWords2Count {
		hotWords2List = append(hotWords2List, KVPair{
			string(k.a) + string(k.b),
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
	allHotSentences = make([][]string, ALL_HOT_LEN_MAX - ALL_HOT_LEN_MIN + 1)
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
			if len(runes) >= ALL_HOT_LEN_MAX - 2 {
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
				allHotSentences[i] = append(allHotSentences[i], s)
			out:
			}
		}
	}
	fmt.Println("仅由不重复的高频字组成的句子")
	for i, c := range allHotSentences {
		fmt.Printf("%d 字：%d 句\n", i + ALL_HOT_LEN_MIN, len(c))
	}

	rand.Seed(1023)
}

// 对称删除模糊匹配算法

// 词典，包含一个 hash 对应的所有句子位置 (文章编号, 句子下标)
var symDelDict map[uint64][]IntPair

// 过滤器，所有错误 hash 的集合
var symDelFilter map[uint64]struct{}

// 初始化词典
func initSymDel() {
	symDelDict = map[uint64][]IntPair{}
	symDelFilter = map[uint64]struct{}{}
	for i, article := range articles {
		for j, s := range article.Content {
			runes := []rune(s)
			hash := uint64(0)
			for _, c := range runes {
				hash += uint64(c) * (uint64(c) + 97)
			}
			// 加入词典
			value := symDelDict[hash]
			if value == nil {
				value = []IntPair{}
			}
			value = append(value, IntPair{i, j})
			symDelDict[hash] = value
			// 加入过滤器
			symDelFilter[hash] = struct{}{}
			for p, c := range runes {
				pval := uint64(c) * (uint64(c) + 97)
				for q := -1; q < 0 && q < p; q++ {
					qval := uint64(0)
					if q != -1 {
						qval = uint64(runes[q]) * (uint64(runes[q]) + 97)
					}
					symDelFilter[hash-pval-qval] = struct{}{}
				}
			}
		}
	}
}

// 检查句子是否在诗词库中
// 返回：(篇目编号, 句子下标)
// 找不到时，返回的第一个值为 -2；若有接近，则为 -1
func checkSentenceInDataset(text []string) (int, int) {
	pivot := 0
	// 找到第一个长度 >= 4 的句子；若无，则取第一句
	for i, s := range text {
		if len([]rune(s)) >= 4 {
			pivot = i
			break
		}
	}

	// 查找数据库
	hash := uint64(0)
	for _, c := range text[pivot] {
		hash += uint64(c) * (uint64(c) + 97)
	}
	if list := symDelDict[hash]; list != nil {
		for _, pos := range list {
			article := &articles[pos.a]
			start := pos.b - pivot
			if start >= 0 && start+len(text) <= len(article.Content) {
				valid := true
				for i, s := range text {
					if article.Content[start+i] != s {
						valid = false
						break
					}
				}
				if valid {
					return pos.a, start
				}
			}
		}
	}

	// 检查是否有接近
	near := true
nearLoop:
	for _, s := range text {
		runes := []rune(s)
		hash := uint64(0)
		for _, c := range runes {
			hash += uint64(c) * (uint64(c) + 97)
		}
		if _, has := symDelFilter[hash]; has {
			continue nearLoop
		}
		for p, c := range runes {
			pval := uint64(c) * (uint64(c) + 97)
			for q := -1; q < p; q++ {
				qval := uint64(0)
				if q != -1 {
					qval = uint64(runes[q]) * (uint64(runes[q]) + 97)
				}
				if _, has := symDelFilter[hash-pval-qval]; has {
					continue nearLoop
				}
			}
		}
		near = false
		break
	}

	if near {
		return -1, -1
	} else {
		return -2, -2
	}
}
