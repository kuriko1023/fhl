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

// 检查一句诗词是否符合规则
func check(sentence string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(sentence, keyword) {
			return false
		}
	}
	return true
}

// 返回一句诗词中的所有高频词，按出现次数降序排序；若无，返回空列表
func getHotWords(sentence string) ([]string, []string) {
	var singleHotWords []string
	var doubleHotWords []string
	for _, str := range hotWords1 {
		if strings.Index(sentence, str) != -1 {
			singleHotWords = append(singleHotWords, str)
		}
	}
	for _, str := range hotWords2 {
		if strings.Index(sentence, str) != -1 {
			doubleHotWords = append(doubleHotWords, str)
		}
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
	for _, i := range randomSample(n, count1 - 1) {
		ret = append(ret, hotWords1[i])
	}
	ret = append(ret, hotWords1[rand.Intn(n) + n])

	// 双字
	n = len(hotWords2) / 2
	for _, i := range randomSample(n, count2 - 1) {
		ret = append(ret, hotWords2[i])
	}
	ret = append(ret, hotWords2[rand.Intn(n) + n])

	return ret
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
	singleWordLimit := sizeRight / sizeLeft + 2
	word2Limit := sizeRight / 5

	for len(right) < sizeRight {
		// 所有已选字的集合
		all := map[string]struct{}{}

		// 先选取 n 个高频字
		left = []string{}
		for _, i := range randomSample(n / 4, leftTopPart) {
			left = append(left, hotWords1[i])
		}
		for _, i := range randomSample(n - n / 4, leftBottomPart) {
			left = append(left, hotWords1[i + n / 4])
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
				if count[j] >= singleWordLimit { continue }
				for i, s := range article.Content {
					if len(s) >= 4 && strings.Contains(s, t) {
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
		//由诗句获得的单字高频词组或双字高频词组为空
		if len(sHotWords) == 0 || len(dHotWords) == 0 {
			continue
		}

		//在双字高频词组长度大于1的情况下， 有0.2的概率贡献两个双字高频词填入题目
		k := rand.Intn(5)
		if k >= 4 && len(dHotWords) > 1 {
			hotWordsList1 = append(hotWordsList2, dHotWords[1])
		}else {
			flag := 0
			//选择第一个与dHotWords[0]无重的单字
			for _, str := range sHotWords {
				if strings.Index(dHotWords[0], str) == -1 {
					hotWordsList1 = append(hotWordsList1, str)
					flag = 1
					goto success
				}
			}
			if flag == 0 {
				continue
			}
		}
		success:
		hotWordsList2 = append(hotWordsList2, dHotWords[0])

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

	type RunePair struct{ a, b rune }
	hotWords1Count := map[rune]int{}
	hotWords2Count := map[RunePair]int{}

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

	rand.Seed(1023)
}
