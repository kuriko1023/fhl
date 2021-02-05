package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
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
	body, err := ioutil.ReadAll(resp.Body)
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
	object := map[string]interface{}{
		"type":        "room_status",
		"host":        GetPlayer(room.Host).Nickname,
		"host_status": "absent",
		"guest":       nil,
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
		object["guest"] = g
	}
	// 向所有玩家的连接发送消息
	for _, p := range room.People {
		p.Channel <- object
	}
}

// 一位玩家选择坐下
func playerSetReady(p *Player) bool {
	room := p.InRoom
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

// 处理玩家客户端发来的消息
func handlePlayerMessage(p *Player, object map[string]interface{}) {
	switch object["type"] {
	case "ready":
		if !playerSetReady(p) {
			// TODO: 错误信息？
			p.Channel <- struct{}{}
		}
	default:
		p.Channel <- struct{}{}
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

	// 写入连接状态
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

	// 建立连接
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	inChannel := make(chan map[string]interface{}, 3)
	outChannel := player.Channel

	go func(c *websocket.Conn, ch chan map[string]interface{}) {
		var object map[string]interface{}
		for {
			if err := c.ReadJSON(&object); err != nil {
				// NOTE: Go 1.16 起使用 net.ErrorClosed
				// https://github.com/golang/go/issues/4373
				if !websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
				) && !strings.Contains(err.Error(), "use of closed network connection") {
					log.Println(err)
				}
				break
			}
			ch <- object
		}
		close(ch)
	}(c, inChannel)

	broadcastRoomStatus(room)

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
		}
	}

	c.Close()

	// 清除状态
	DataMutex.Lock()
	room.People = removeElement(room.People, player)
	if player.Id == room.Host {
		room.HostReady = false
	} else if player.Id == room.Guest {
		room.Guest = ""
	}
	player.InRoom = nil
	player.Channel = nil
	DataMutex.Unlock()

	// 即使 inChannel 尚未关闭，它也将由
	// goroutine 在 c.Close() 触发的错误之后关闭
	// 而 outChannel 不能被外界关闭
	close(outChannel)

	broadcastRoomStatus(room)

	log.Println("connection closed")
}

var testCounter = 0

func testHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("test.html")
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

func SetUpHttp() {
	http.HandleFunc("/channel/", channelHandler)
	if Config.Debug {
		http.HandleFunc("/test", testHandler)
	}

	port := Config.Port
	log.Printf("Listening on http://localhost:%d/\n", port)
	if Config.Debug {
		log.Printf("Visit http://localhost:%d/test for testing\n", port)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
