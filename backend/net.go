package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	cumulativeTimer = 60000
	turnTimer       = 60000
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 登录
// 传入微信登录 wx.login() 获得的 code
// 返回 OpenID
func login(code string) string {
	// 调试用，若 code 以 "!" 开头，则 OpenID 等于 code 去掉 "!"
	if Config.Debug && len(code) >= 2 && code[0] == '!' {
		return code[1:]
	}

	resp, err := http.Get(
		"https://api.weixin.qq.com/sns/jscode2session?appid=" + Config.AppID +
			"&secret=" + Config.AppSecret +
			"&js_code=" + code +
			"&grant_type=authorization_code",
	)
	if err != nil {
		log.Println(err)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	// 解析 JSON 数据
	var object struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
	}
	err = json.Unmarshal(body, &object)
	if err != nil {
		log.Println(err)
		return ""
	}
	if object.OpenID == "" || object.SessionKey == "" {
		log.Println("Invalid response", string(body))
		return ""
	}

	return object.OpenID
}

func retrieveAvatar(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") || strings.IndexByte(contentType, ' ') != -1 {
		return []byte{}
	}
	body, err := io.ReadAll(http.MaxBytesReader(nil, resp.Body, 64*1024))
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	// TODO: Store the MIME type
	return body
}

// 移除 slice 中的一个元素
func removeElement(s []*Player, p *Player) []*Player {
	for i, e := range s {
		if e == p {
			s[i] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

// 向所有房间中的玩家更新房间状态
func broadcastRoomStatus(room *Room) {
	p := GetPlayer(room.Host)
	object := map[string]interface{}{
		"type":         "room_status",
		"host":         p.Nickname,
		"host_avatar":  p.Id + "/" + strconv.FormatInt(p.AvatarUpdated, 36),
		"host_status":  "absent",
		"guest":        nil,
		"guest_avatar": nil,
	}
	if room.HostReady {
		object["host_status"] = "ready"
	} else {
		hostPresent := false
		for _, p := range room.People {
			if p.Id == room.Host {
				hostPresent = true
				break
			}
		}
		if hostPresent {
			object["host_status"] = "present"
		}
	}
	if g := room.Guest; g != "" {
		p := GetPlayer(g)
		object["guest"] = p.Nickname
		object["guest_avatar"] = p.Id + "/" + strconv.FormatInt(p.AvatarUpdated, 36)
	}
	// 向所有玩家的连接发送消息
	for _, p := range room.People {
		if p.Channel != nil {
			p.Channel <- object
		}
	}
}

// 一位玩家选择坐下
func playerReady(p *Player) bool {
	room := p.InRoom
	if room.State != "" {
		return false
	}
	if room.Host == p.Id {
		room.HostReady = true
	} else {
		if room.Guest != "" && room.Guest != p.Id {
			// 已经有人了 T-T
			return false
		}
		room.Guest = p.Id
	}
	broadcastRoomStatus(room)
	return true
}

func roomGenerateStatus(room *Room) map[string]interface{} {
	subjectRepr := interface{}(nil)
	if subject := room.Subject; subject != nil {
		subjectRepr = subject.Dump()
	}
	size, _ := strconv.Atoi(room.Mode[1:])
	return map[string]interface{}{
		"type":    "generated",
		"mode":    room.Mode[0:1],
		"size":    size,
		"subject": subjectRepr,
	}
}

func roomHistoryStrings(room *Room) []string {
	history := []string{}
	for _, a := range room.History {
		history = append(history, a.Dump())
	}
	return history
}

func roomGameStatus(room *Room) map[string]interface{} {
	hostTimer := room.HostTimer
	guestTimer := room.GuestTimer
	curTurnTimer := turnTimer - int(nowMilliseconds()-room.LastMoveAt)
	if curTurnTimer < 0 {
		if room.CurMoveSide == SideHost {
			hostTimer += curTurnTimer
		} else {
			guestTimer += curTurnTimer
		}
		curTurnTimer = 0
	}
	return map[string]interface{}{
		"type":        "game_status",
		"mode":        room.Mode,
		"subject":     room.Subject.Dump(),
		"history":     roomHistoryStrings(room),
		"host_timer":  hostTimer,
		"guest_timer": guestTimer,
		"turn_timer":  curTurnTimer,
	}
}

func bicastGameDelta(room *Room, change interface{}) {
	object := map[string]interface{}{
		"type":        "game_update",
		"text":        room.History[len(room.History)-1].Dump(),
		"update":      change,
		"host_timer":  room.HostTimer,
		"guest_timer": room.GuestTimer,
	}
	if Players[room.Host].Channel != nil {
		Players[room.Host].Channel <- object
	}
	if room.Guest != "" && Players[room.Guest].Channel != nil {
		Players[room.Guest].Channel <- object
	}
}

func bicastGameEnd(room *Room, winner Side) {
	var winnerVal int
	switch winner {
	case SideHost:
		winnerVal = 1
	case SideGuest:
		winnerVal = -1
	case SideNone:
		winnerVal = 0
	}
	if Players[room.Host].Channel != nil {
		Players[room.Host].Channel <- map[string]interface{}{
			"type":    "end_status",
			"winner":  winnerVal,
			"mode":    room.Mode,
			"subject": room.Subject.Dump(),
			"history": roomHistoryStrings(room),
		}
	}
	if room.Guest != "" && Players[room.Guest].Channel != nil {
		Players[room.Guest].Channel <- map[string]interface{}{
			"type":    "end_status",
			"winner":  -winnerVal,
			"mode":    room.Mode,
			"subject": room.Subject.Dump(),
			"history": roomHistoryStrings(room),
		}
	}
}

func nowMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func resetRoomState(room *Room) {
	room.State = ""
	room.HostReady = false
	room.Guest = ""
	room.Subject = nil
}

func errorMsg(s string) map[string]string {
	return map[string]string{"error": s}
}

func parseMode(x interface{}) string {
	if x == "A" || x == "B" || x == "C" || x == "D" {
		return x.(string)
	} else {
		return ""
	}
}

func parseInt(x interface{}) int {
	f, ok := x.(float64)
	if !ok {
		return -1
	}
	return int(math.Round(f))
}

// 处理玩家客户端发来的消息
// NOTE: 大部分业务逻辑在此处实现
func handlePlayerMessage(p *Player, object map[string]interface{}) {
	if p.InRoom != nil {
		m := p.InRoom.Mutex
		m.Lock()
		defer m.Unlock()
	}

	defer func() {
		if e := recover(); e != nil {
			if s, ok := e.(string); ok {
				if p.Channel != nil {
					p.Channel <- errorMsg(s)
				}
			}
		}
	}()

	switch object["type"] {
	case "profile":
		nickname, ok1 := object["nickname"].(string)
		avatar, ok2 := object["avatar"].(string)
		// TODO: 用 session key 验证签名
		if !ok1 || !ok2 {
			panic("Incorrect format")
		}
		p.Nickname = nickname
		p.Avatar = retrieveAvatar(avatar)
		p.AvatarUpdated = time.Now().Unix()
		if err := p.Save(); err != nil {
			fmt.Println(err)
		}
		if p.InRoom != nil {
			broadcastRoomStatus(p.InRoom)
		}

	case "ready":
		if !playerReady(p) {
			panic("Already occupied")
		}

	case "start_generate":
		if p.InRoom == nil || p.InRoom.Host != p.Id {
			panic("Must be host")
		}
		if p.InRoom.State != "" || !p.InRoom.HostReady || p.InRoom.Guest == "" {
			panic("Room should be idle with two ready players")
		}
		p.InRoom.State = "gen"
		p.InRoom.Mode = "A"
		if p.InRoom.Guest != "" && Players[p.InRoom.Guest].Channel != nil {
			Players[p.InRoom.Guest].Channel <- roomGenerateStatus(p.InRoom)
		}

	case "set_mode":
		fallthrough
	case "generate":
		isGenerate := (object["type"] == "generate")
		if p.InRoom == nil || p.InRoom.Host != p.Id {
			panic("Must be host")
		}
		if p.InRoom.State != "gen" {
			panic("Room should be in generation phase")
		}
		mode := parseMode(object["mode"])
		size := parseInt(object["size"])
		if mode == "" || size == -1 {
			panic("Incorrect format")
		}

		// 检查 size
		switch mode {
		case "A":
			if size != 0 {
				panic("Incorrect size")
			}
		case "B":
			if size < 5 || size > 9 {
				panic("Incorrect size")
			}
		case "C":
			if size != 1 && size != 3 {
				panic("Incorrect size")
			}
		case "D":
			if size < 5 || size > 10 {
				panic("Incorrect size")
			}
		}
		p.InRoom.Mode = mode + strconv.Itoa(size)

		// 生成题目
		if isGenerate {
			var subject Subject
			switch mode {
			case "A":
				var word string
				if rand.Intn(4) == 0 {
					word = generateA(0, 1)[0]
				} else {
					word = generateA(1, 0)[0]
				}
				subject = &SubjectA{Word: word}
			case "B":
				words := generateB(size)
				subject = &SubjectB{Words: []rune(words), CurIndex: 0}
			case "C":
				left, right := generateC(size, 7+3*size)
				subject = &SubjectC{
					WordsLeft:  left,
					WordsRight: right,
					UsedRight:  make([]bool, 10),
				}
			case "D":
				left, right := generateD(size)
				subject = &SubjectD{
					WordsLeft:  left,
					WordsRight: right,
					UsedLeft:   make([]bool, size),
					UsedRight:  make([]bool, size),
				}
			}
			p.InRoom.Subject = subject
		} else {
			p.InRoom.Subject = nil
		}

		object := roomGenerateStatus(p.InRoom)
		if p.InRoom.Guest != "" && Players[p.InRoom.Guest].Channel != nil {
			Players[p.InRoom.Guest].Channel <- object
		}
		if isGenerate {
			p.Channel <- object
		}

	case "start_game":
		if p.InRoom == nil || p.InRoom.Host != p.Id {
			panic("Must be host")
		}
		if p.InRoom.State != "gen" {
			panic("Room should be in generation phase")
		}
		if p.InRoom.Subject == nil {
			panic("No subject generated")
		}

		p.InRoom.State = "game"
		p.InRoom.Mode = p.InRoom.Mode[0:1]
		p.InRoom.History = []CorrectAnswer{}
		p.InRoom.HistorySet = map[string]struct{}{}
		p.InRoom.LastMoveAt = nowMilliseconds()
		p.InRoom.CurMoveSide = SideHost
		p.InRoom.HostTimer = cumulativeTimer
		p.InRoom.GuestTimer = cumulativeTimer

		object := roomGameStatus(p.InRoom)
		if Players[p.InRoom.Host].Channel != nil {
			Players[p.InRoom.Host].Channel <- object
		}
		if p.InRoom.Guest != "" && Players[p.InRoom.Guest].Channel != nil {
			Players[p.InRoom.Guest].Channel <- object
		}

		p.InRoom.TimerStopSignal = make(chan struct{})
		go func(room *Room) {
			m := room.Mutex
			ch := room.TimerStopSignal
			for {
				select {
				case <-ch:
					return
				case <-time.After(time.Second / 2):
					m.Lock()
					var curTimer int
					var oppTimer int
					if room.CurMoveSide == SideHost {
						curTimer = room.HostTimer
						oppTimer = room.GuestTimer
					} else {
						curTimer = room.GuestTimer
						oppTimer = room.HostTimer
					}
					if oppTimer < 0 {
						// 小概率事件，对方上次提交答案时已经超时
						bicastGameEnd(room, room.CurMoveSide)
						resetRoomState(room)
						room.TimerStopSignal = nil
						m.Unlock()
						close(ch)
						return
					}
					timeUsed := int(nowMilliseconds()-room.LastMoveAt) - turnTimer
					if timeUsed < 0 {
						timeUsed = 0
					}
					if curTimer < timeUsed {
						bicastGameEnd(room, 1-room.CurMoveSide)
						resetRoomState(room)
						room.TimerStopSignal = nil
						m.Unlock()
						close(ch)
						return
					}
					m.Unlock()
				}
			}
		}(p.InRoom)

	case "answer":
		if p.InRoom.State != "game" {
			panic("Room should be in game phase")
		}
		playerSide := SideGuest
		if p.Id == p.InRoom.Host {
			playerSide = SideHost
		}
		if p.InRoom.CurMoveSide != playerSide {
			panic("Not your move")
		}
		text, ok := object["text"].(string)
		if !ok {
			panic("Incorrect format")
		}

		incorrectReason := ""

		texts := strings.Split(text, " ")
		correct, articleIdx, sentenceIdx := lookupText(texts)
		if !correct {
			if articleIdx != -1 {
				incorrectReason = "没背熟"
			} else {
				incorrectReason = "大文豪"
			}
			println(articleIdx, sentenceIdx)
		}

		if incorrectReason == "" {
			for _, s := range texts {
				if _, ok := p.InRoom.HistorySet[s]; ok {
					incorrectReason = "复读机"
					break
				}
			}
		}

		// 调试模式下叹号开头的可以免受语料库与历史限制
		if Config.Debug && strings.HasPrefix(text, "!") {
			incorrectReason = ""
		}

		var kws []IntPair
		var change interface{}
		if incorrectReason == "" {
			kws, change = p.InRoom.Subject.Answer(text, playerSide)
			if kws == nil {
				incorrectReason = "不审题"
			}
		}

		if incorrectReason != "" {
			if p.Channel != nil {
				p.Channel <- map[string]string{
					"type":   "invalid_answer",
					"reason": incorrectReason,
				}
			}
			break
		}

		p.InRoom.History = append(p.InRoom.History, CorrectAnswer{text, kws})
		for _, s := range texts {
			p.InRoom.HistorySet[s] = struct{}{}
		}

		timeNow := nowMilliseconds()
		timeUsed := int(timeNow-p.InRoom.LastMoveAt) - turnTimer
		if timeUsed < 0 {
			timeUsed = 0
		}
		if p.InRoom.CurMoveSide == SideHost {
			p.InRoom.HostTimer -= timeUsed
		} else {
			p.InRoom.GuestTimer -= timeUsed
		}

		p.InRoom.LastMoveAt = timeNow
		p.InRoom.CurMoveSide = 1 - p.InRoom.CurMoveSide

		// 游戏是否已经完成（用完所有的字词）
		if p.InRoom.Subject.End() {
			p.InRoom.TimerStopSignal <- struct{}{}
			p.InRoom.TimerStopSignal = nil
			bicastGameEnd(p.InRoom, SideNone)
			resetRoomState(p.InRoom)
		} else {
			bicastGameDelta(p.InRoom, change)
		}

	default:
		panic("Unknown type")
	}
}

func channelHandler(w http.ResponseWriter, r *http.Request) {
	cmpns := strings.SplitN(r.URL.Path[len("/channel/"):], "/", 2)
	if len(cmpns) != 2 {
		// Bad Request，URL 不合法
		w.WriteHeader(400)
		return
	}

	roomId := cmpns[0]
	id := login(cmpns[1])

	// 确认登录信息有效
	if id == "" {
		// Unauthorized，登录 code 不合法
		w.WriteHeader(401)
		return
	}

	// 查找/加入玩家数据库
	player := GetPlayer(id)

	// 确认房间有效
	if roomId == "my" {
		// 进入自己的房间
		roomId = player.Id
	}
	room := Rooms[roomId]
	if room == nil {
		// Not Found，没有对应的房间
		w.WriteHeader(404)
		return
	}

	// 更新全局存储的连接状态
	DataMutex.Lock()
	if player.Channel != nil {
		player.Channel <- nil
		// 复制一份
		player = &Player{
			Id:       player.Id,
			Nickname: player.Nickname,
			Avatar:   player.Avatar,
		}
		Players[player.Id] = player
	}
	player.Channel = make(chan interface{}, 3)
	player.InRoom = room
	room.People = append(room.People, player)
	DataMutex.Unlock()

	// 建立 WebSocket 连接
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置读限制
	c.SetReadLimit(4096)
	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	// inChannel: 客户端发来的消息，以 string 到 interface{} 的 map 表示
	// outChannel: 要发送至客户端的消息，即 player.Channel
	inChannel := make(chan map[string]interface{}, 3)
	outChannel := player.Channel

	// Goroutine，不断从 WebSocket 连接读入 JSON 并发送至 inChannel
	go func(c *websocket.Conn, ch chan map[string]interface{}) {
		var object map[string]interface{}
		for {
			if err := c.ReadJSON(&object); err != nil {
				if !websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) && !errors.Is(err, net.ErrClosed) {
					log.Println(err)
				}
				if _, ok := err.(*json.SyntaxError); ok {
					continue
				}
				break
			}
			ch <- object
		}
		close(ch)
	}(c, inChannel)

	if room.State == "" {
		broadcastRoomStatus(room)
	} else if room.State == "gen" {
		outChannel <- roomGenerateStatus(room)
	} else if room.State == "game" {
		outChannel <- roomGameStatus(room)
	}

	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

messageLoop:
	for inChannel != nil && outChannel != nil {
		select {
		case object, ok := <-inChannel:
			if !ok {
				inChannel = nil
				break messageLoop
			}
			handlePlayerMessage(player, object)

		case object := <-outChannel:
			if object == nil {
				break messageLoop
			}
			if err := c.WriteJSON(object); err != nil {
				log.Println(err)
				break messageLoop
			}

		case <-pingTicker.C:
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				break messageLoop
			}
		}
	}

	c.Close()

	// 清除全局存储的状态
	DataMutex.Lock()
	room.People = removeElement(room.People, player)
	if room.State == "" {
		// 只有房间处于等待状态时，才认为是退出了房间
		if player.Id == room.Host {
			room.HostReady = false
		} else if player.Id == room.Guest {
			room.Guest = ""
		}
	}
	player.InRoom = nil
	player.Channel = nil
	DataMutex.Unlock()

	// 即使 inChannel 尚未关闭，它也将由
	// goroutine 在 c.Close() 触发的错误之后关闭
	// 而 outChannel 不能被外界关闭
	close(outChannel)

	if room.State == "" {
		broadcastRoomStatus(room)
	}

	log.Println("connection closed")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	id := login(r.URL.Path[len("/profile/"):])

	obj := map[string]interface{}{
		"id":       id,
		"nickname": nil,
	}
	if p := Players[id]; p != nil {
		obj["nickname"] = p.Nickname
	}
	enc := json.NewEncoder(w)
	enc.Encode(obj)
}

var defaultAvatar []byte

func init() {
	var err error
	defaultAvatar, err = os.ReadFile("../frontend/src/static_remote/grey_avatar_132.jpg")
	if err != nil {
		log.Fatalln(err)
	}
}

func avatarHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.SplitN(r.URL.Path[len("/avatar/"):], "/", 2)[0]
	if id == "" {
		w.Write(defaultAvatar)
		return
	}
	if p := Players[id]; p != nil {
		w.Write(p.Avatar)
	} else {
		w.WriteHeader(404)
	}
}

var testCounter = 0

func testHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("test.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	s := string(content)
	var host string
	var me string
	if testCounter%2 == 0 {
		host = "my"
		me = "kuriko1023"
	} else {
		host = "kuriko1023"
		me = "PiscesOvO"
	}
	testCounter++
	s = strings.Replace(s, "~ host ~", host, 1)
	s = strings.Replace(s, "~ me ~", me, 1)
	w.Write([]byte(s))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	DataMutex.Lock()
	for _, room := range Rooms {
		if room.TimerStopSignal != nil {
			room.TimerStopSignal <- struct{}{}
			room.TimerStopSignal = nil
		}
		resetRoomState(room)
	}
	DataMutex.Unlock()
}

func SetUpHttp() {
	http.HandleFunc("/channel/", channelHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/avatar/", avatarHandler)
	if Config.Debug {
		http.HandleFunc("/test", testHandler)
		http.HandleFunc("/reset", resetHandler)
	}
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("../frontend/src/static_remote"))))

	port := Config.Port
	log.Printf("Listening on http://localhost:%d/\n", port)
	if Config.Debug {
		log.Printf("Visit http://localhost:%d/test for testing\n", port)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
