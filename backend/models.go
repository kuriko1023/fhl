package main

import (
	"strconv"
	"strings"
)

type Player struct {
	Id       string
	InRoom   string
	Nickname string
}

// 表示游戏中的一方，房主或客人
type Side int

const (
	SideHost  Side = 0
	SideGuest Side = 1
)

type Room struct {
	Host       string  // 房主 id
	Guest      string  // 客人 id
	HostReady  bool    // 房主是否准备
	GuestReady bool    // 客人是否准备
	Mode       string  // 游戏模式，空字符串或 "A" "B" "C" "D" 之一
	Subject    Subject // 游戏题目与进度（详细见下）
	// 之前提交的所有文本，偶数下标对应房主，奇数下标对应客人
	// 若一次提交包含多句（以标点分隔的小段），则用斜杠“/”分隔
	History    []string
	CurMove    Side // 当前正答题的一方
	HostTimer  int  // 房主的剩余时间
	GuestTimer int  // 客人的剩余时间
}

// 游戏题目与进度的接口，下面的 SubjectA/B/C/D 都会实现之。
// 包含任意形式的数据结构，但是需要支持与字符串的互转，
// 以及「尝试提交一个句子，若正确则更新状态」的操作。
type Subject interface {
	Parse(string) // 从一个字符串解析出数据
	Dump() string // 将数据表示为一个字符串
	// 尝试用一段文本作答，若与题目匹配则更新数据结构并返回 true，否则返回 false
	// 第一个参数是提交的文本内容，包含多句时用斜杠分隔
	// 第二个参数为当前答题的一方
	Answer(string, Side) (bool, interface{})
}

// 普通飞花 题目与进度
type SubjectA struct {
	Word string
}

func (s *SubjectA) Parse(str string) {
	s.Word = str
}
func (s *SubjectA) Dump() string {
	return s.Word
}
func (s *SubjectA) Answer(str string, from Side) (bool, interface{}) {
	return strings.Contains(str, s.Word), nil
}

// 多字飞花 题目与进度
type SubjectB struct {
	Words    []rune
	CurIndex int
}

func (s *SubjectB) Parse(str string) {
	// 例：春花秋月何时了/3
	fields := strings.SplitN(str, "/", 2)
	s.Words = []rune(fields[0])
	s.CurIndex, _ = strconv.Atoi(fields[1])
}
func (s *SubjectB) Dump() string {
	return string(s.Words) + "/" + strconv.Itoa(s.CurIndex)
}
func (s *SubjectB) Answer(str string, from Side) (bool, interface{}) {
	if strings.ContainsRune(str, s.Words[s.CurIndex]) {
		if from == SideGuest {
		}
		return true, nil
	} else {
		return false, nil
	}
}

// 超级飞花 题目与进度
type SubjectC struct {
	WordsLeft  []string
	WordsRight []string
	UsedRight  []bool
}

func (s *SubjectC) Parse(str string) {
	// 例：古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011
	fields := strings.SplitN(str, "/", 3)
	s.WordsLeft = strings.Split(fields[0], " ")
	s.WordsRight = strings.Split(fields[1], " ")
	s.UsedRight = make([]bool, len(s.WordsRight))
	for i := range s.UsedRight {
		s.UsedRight[i] = (fields[2][i] == '1')
	}
}
func (s *SubjectC) Dump() string {
	used := []rune{}
	for _, b := range s.UsedRight {
		if b {
			used = append(used, '1')
		} else {
			used = append(used, '0')
		}
	}
	return strings.Join(s.WordsLeft, " ") + "/" +
		strings.Join(s.WordsRight, " ") + "/" +
		string(used)
}
func (s *SubjectC) Answer(str string, from Side) (bool, interface{}) {
	indexLeft, indexRight := -1, -1
	for i, w := range s.WordsLeft {
		if strings.Contains(str, w) {
			indexLeft = i
			break
		}
	}
	for i, w := range s.WordsRight {
		if !s.UsedRight[i] && strings.Contains(str, w) {
			indexRight = i
			break
		}
	}
	if indexLeft == -1 || indexRight == -1 {
		return false, nil
	}
	s.UsedRight[indexRight] = true
	return true, [2]int{indexLeft, indexRight}
}

// 谜之飞花 题目与进度
type SubjectD struct {
	WordsLeft  []string
	WordsRight []string
	UsedLeft   []bool
	UsedRight  []bool
}

func (s *SubjectD) Parse(str string) {
}
func (s *SubjectD) Dump() string {
	return ""
}
func (s *SubjectD) Answer(str string, from Side) (bool, interface{}) {
	return false, nil
}
