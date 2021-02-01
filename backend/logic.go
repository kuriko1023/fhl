package main

import (
	"bufio"
	"fmt"
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

// 所有高频词组成的列表，包含单字和双字
// 键表示高频词，值表示该词出现的频数
var hotWords map[string]int

// 检查一句诗词是否符合规则
func check(sentence string, keywords []string) bool {
	return false
}

// 返回一句诗词中的所有高频词，按出现次数降序排序；若无，返回空列表
func getHotWords(sentence string) []string {
	return nil
}

// 生成普通飞花题目备选列表，返回 count 个不重复的单字/词
func generateA(count int) []string {
	return nil
}

// 生成超级飞花题目，返回长度为 n 和 m 的字符串列表
func generateB() ([]string, []string) {
	return nil, nil
}

// 生成 2 个长度为 n 的关键词组作为谜之飞花题目
func generateC() ([]string, []string) {
	return nil, nil
}

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
	file, err := os.Open("dataset.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	type RunePair struct{ a, b rune }
	hotWords1 := map[rune]int{}
	hotWords2 := map[RunePair]int{}

	i := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		// 随机抽取十分之一
		i++
		if i%10 != 0 {
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
					hotWords1[c] += weight
					if i < len(s)-1 {
						hotWords2[RunePair{c, s[i+1]}] += weight
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
	for k, v := range hotWords1 {
		hotWords1List = append(hotWords1List, KVPair{string(k), v})
	}
	for k, v := range hotWords2 {
		hotWords2List = append(hotWords2List, KVPair{
			string(k.a) + string(k.b),
			v,
		})
	}
	sort.Sort(hotWords1List)
	sort.Sort(hotWords2List)

	fmt.Println("高频单字，每行五十个")
	for i := 0; i < 500; i++ {
		fmt.Print(hotWords1List[i].string)
		if (i+1)%50 == 0 {
			fmt.Println()
		}
	}
	fmt.Println("高频双字，每行二十个")
	for i := 0; i < 200; i++ {
		fmt.Print(hotWords2List[i].string, " ")
		if (i+1)%20 == 0 {
			fmt.Println()
		}
	}
}
