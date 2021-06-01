package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"database/sql"
)

type Player struct {
	Id       string
	Nickname string
	Avatar   []byte

	// 头像上次更新的时刻
	AvatarUpdated int64

	// 当前所处的房间
	InRoom *Room

	// 传送消息的信道
	// 往这里发送的值会在 JSON 编码后经 WebSocket 连接传至客户端
	Channel chan interface{}
}

// 表示游戏中的一方，房主或客人
type Side int

const (
	SideHost  Side = 0
	SideGuest Side = 1
	SideNone  Side = -1
)

type IntPair struct {
	a, b int
}

type CorrectAnswer struct {
	string
	Keywords []IntPair
}

type Room struct {
	Host      string  // 房主 id
	HostReady bool    // 房主是否坐下
	Guest     string  // 已经坐下的客人 id
	State     string  // 游戏状态，空字符串表示等待，gen 表示正在选题，game 表示游戏中
	Mode      string  // "A" "B" "C" "D" 之一表示游戏玩法。在 gen 状态下，后接空格与 size
	Subject   Subject // 游戏题目与进度（详细见下）
	// 之前提交的所有文本，偶数下标对应房主，奇数下标对应客人
	// 若一次提交包含多句（以标点分隔的小段），则用斜杠“/”分隔
	History []CorrectAnswer
	// 之前提交的所有文本中以标点分隔的小段集合
	HistorySet  map[string]struct{}
	LastMoveAt  int64 // 上一次提交答案的时间戳
	CurMoveSide Side  // 当前正答题的一方
	HostTimer   int   // 房主的剩余时间
	GuestTimer  int   // 客人的剩余时间

	People []*Player // 建立了此房间的 WebSocket 连接的人

	TimerStopSignal chan struct{}
	Mutex           *sync.Mutex
}

func (a CorrectAnswer) Dump() string {
	b := strings.Builder{}
	b.WriteString(a.string)
	for _, k := range a.Keywords {
		b.WriteRune('/')
		for i := 0; i < k.b; i++ {
			b.WriteString(strconv.FormatInt(int64(k.a+i), 36))
		}
	}
	return b.String()
}

// 游戏题目与进度的接口，下面的 SubjectA/B/C/D 都会实现之。
// 包含任意形式的数据结构，但是需要支持与字符串的互转，
// 以及「尝试提交一个句子，若正确则更新状态」的操作。
type Subject interface {
	Parse(string) // 从一个字符串解析出数据
	Dump() string // 将数据表示为一个字符串
	// 尝试用一段文本作答，若与题目匹配则更新数据结构
	// 第一个参数是提交的文本内容，包含多句时用斜杠分隔
	// 第二个参数为当前答题的一方
	// 第一个返回值是关键词的下标集合，若答案不合法则为 nil
	// 第二个返回值是一个字符串，表示变化量
	Answer(string, Side) ([]IntPair, string)
	End() bool // 游戏是否结束
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
func (s *SubjectA) Answer(str string, from Side) ([]IntPair, string) {
	// 第二个返回值：""
	p := strings.Index(str, s.Word)
	if p != -1 {
		return []IntPair{IntPair{runes(str[:p]), runes(s.Word)}}, ""
	} else {
		return nil, ""
	}
}
func (s *SubjectA) End() bool {
	return false
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
func (s *SubjectB) Answer(str string, from Side) ([]IntPair, string) {
	// 第二个返回值：下一位轮到的玩家要飞的字的下标，若游戏结束则为 -1
	p := strings.IndexRune(str, s.Words[s.CurIndex])
	if p != -1 {
		if from == SideGuest {
			s.CurIndex++
			if s.CurIndex == len(s.Words) {
				s.CurIndex = -1
			}
		}
		return []IntPair{IntPair{runes(str[:p]), 1}}, strconv.Itoa(s.CurIndex)
	} else {
		return nil, ""
	}
}
func (s *SubjectB) End() bool {
	return (s.CurIndex == -1)
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
func (s *SubjectC) Answer(str string, from Side) ([]IntPair, string) {
	// 第二个返回值：右侧被匹配的关键词下标
	indexLeft, indexRight := -1, -1
	ps := make([]IntPair, 2)
	for i, w := range s.WordsLeft {
		p := strings.Index(str, w)
		if p != -1 {
			indexLeft = i
			ps[0] = IntPair{runes(str[:p]), runes(w)}
			break
		}
	}
	for i, w := range s.WordsRight {
		p := strings.Index(str, w)
		if !s.UsedRight[i] && p != -1 {
			indexRight = i
			ps[1] = IntPair{runes(str[:p]), runes(w)}
			break
		}
	}
	if indexLeft == -1 || indexRight == -1 {
		return nil, ""
	}
	s.UsedRight[indexRight] = true
	return ps, strconv.Itoa(indexRight)
}
func (s *SubjectC) End() bool {
	for _, u := range s.UsedRight {
		if !u {
			return false
		}
	}
	return true
}

// 谜之飞花 题目与进度
type SubjectD struct {
	WordsLeft  []string
	WordsRight []string
	UsedLeft   []bool
	UsedRight  []bool
}

func (s *SubjectD) Parse(str string) {
	// 例：万 书 今 凉 得 来 柳 欲/一片 丝 如此 孤 庭 细 舟 觉/00000000/00000000
	fields := strings.SplitN(str, "/", 4)
	s.WordsLeft = strings.Split(fields[0], " ")
	s.WordsRight = strings.Split(fields[1], " ")
	s.UsedLeft = make([]bool, len(s.WordsLeft))
	s.UsedRight = make([]bool, len(s.WordsRight))
	for i := range s.UsedLeft {
		s.UsedLeft[i] = (fields[2][i] == '1')
	}
	for i := range s.UsedRight {
		s.UsedRight[i] = (fields[3][i] == '1')
	}
}
func (s *SubjectD) Dump() string {
	used := []rune{}
	for _, b := range s.UsedLeft {
		if b {
			used = append(used, '1')
		} else {
			used = append(used, '0')
		}
	}
	used = append(used, '/')
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
func (s *SubjectD) Answer(str string, from Side) ([]IntPair, string) {
	// 第二个返回值："a,b"，左右侧被匹配的关键词下标
	indexLeft, indexRight := -1, -1
	ps := make([]IntPair, 2)
	for i, w := range s.WordsLeft {
		p := strings.Index(str, w)
		if !s.UsedLeft[i] && p != -1 {
			indexLeft = i
			ps[0] = IntPair{runes(str[:p]), runes(w)}
			break
		}
	}
	for i, w := range s.WordsRight {
		p := strings.Index(str, w)
		if !s.UsedRight[i] && p != -1 {
			indexRight = i
			ps[1] = IntPair{runes(str[:p]), runes(w)}
			break
		}
	}
	if indexLeft == -1 || indexRight == -1 {
		return nil, ""
	}
	s.UsedLeft[indexLeft] = true
	s.UsedRight[indexRight] = true
	return ps, strconv.Itoa(indexLeft) + "," + strconv.Itoa(indexRight)
}
func (s *SubjectD) End() bool {
	for _, u := range s.UsedLeft {
		if !u {
			return false
		}
	}
	for _, u := range s.UsedRight {
		if !u {
			return false
		}
	}
	return true
}

// 计算一个字符串中的 Unicode 字符数
func runes(s string) int {
	return len([]rune(s))
}

// 数据库

// 访问两个全局 map 的互斥锁
var DataMutex = &sync.Mutex{}

// 所有玩家，键为玩家 ID
var Players map[string]*Player

// 所有房间，键为房主 ID
var Rooms map[string]*Room

// 与数据库交互的逻辑

func SetUpDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "fhl.db")
	if err != nil {
		return nil, err
	}

	// 创建表
	cmd := "CREATE TABLE IF NOT EXISTS players" +
		"(id TEXT UNIQUE PRIMARY KEY, nickname TEXT, avatar BLOB)"
	if _, err := db.Exec(cmd); err != nil {
		db.Close()
		return nil, err
	}

	// 清空全局数据
	Players = map[string]*Player{}
	Rooms = map[string]*Room{}

	// 读取玩家信息
	now := time.Now().Unix()
	rows, err := db.Query("SELECT * FROM players")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := Player{AvatarUpdated: now}
		err := rows.Scan(&p.Id, &p.Nickname, &p.Avatar)
		if err != nil {
			return nil, err
		}
		Players[p.Id] = &p
		Rooms[p.Id] = &Room{Host: p.Id, Mutex: &sync.Mutex{}}
	}

	return db, nil
}

func (p *Player) Save() error {
	Players[p.Id] = p
	if Rooms[p.Id] == nil {
		Rooms[p.Id] = &Room{Host: p.Id, Mutex: &sync.Mutex{}}
	}
	_, err := db.Exec("INSERT INTO players(id, nickname, avatar) "+
		"VALUES($1, $2, $3) "+
		"ON CONFLICT(id) DO UPDATE SET "+
		"nickname = excluded.nickname, "+
		"avatar = excluded.avatar",
		p.Id, p.Nickname, p.Avatar)
	return err
}

func GetPlayer(id string) *Player {
	if p := Players[id]; p != nil {
		return p
	}
	p := &Player{
		Id:            id,
		Nickname:      "猫猫" + id,
		Avatar:        []byte{},
		AvatarUpdated: time.Now().Unix(),
	}
	p.Save()
	return p
}
