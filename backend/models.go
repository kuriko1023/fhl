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
	CurMove    int // 当前正答题的一方，0 表示房主，1 表示客人
	HostTimer  int // 房主的剩余时间
	GuestTimer int // 客人的剩余时间
}

// 游戏题目与进度的接口，下面的 SubjectA/B/C/D 都会实现之。
// 包含任意形式的数据结构，但是需要支持与字符串的互转，
// 以及「尝试提交一个句子，若正确则更新状态」的操作。
type Subject interface {
	Parse(string) // 从一个字符串解析出数据
	Dump() string // 将数据表示为一个字符串
	// 尝试提交一段文本，若与题目匹配则更新数据结构并返回 true，否则返回 false
	// 第一个参数是提交的文本内容，包含多句时用斜杠分隔
	// 第二个参数为 false 表示提交的是房主，true 表示是客人
	Answer(string, bool) bool
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
func (s *SubjectA) Answer(str string, from bool) bool {
	return check(str, []string{s.Word})
}

// 多字飞花 题目与进度
type SubjectB struct {
	Words    []rune
	CurIndex int
}

func (s *SubjectB) Parse(str string) {
	// 例：春花秋月何时了 3
	fields := strings.SplitN(str, " ", 2)
	s.Words = []rune(fields[0])
	s.CurIndex, _ = strconv.Atoi(fields[1])
}
func (s *SubjectB) Dump() string {
	return string(s.Words) + " " + strconv.Itoa(s.CurIndex)
}
func (s *SubjectB) Answer(str string, from bool) bool {
	return false
}

// 超级飞花 题目与进度
type SubjectC struct {
	WordsLeft  []string
	WordsRight []string
	UsedLeft   []bool
}

func (s *SubjectC) Parse(str string) {
}
func (s *SubjectC) Dump() string {
	return ""
}
func (s *SubjectC) Answer(str string, from bool) bool {
	return false
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
func (s *SubjectD) Answer(str string, from bool) bool {
	return false
}
