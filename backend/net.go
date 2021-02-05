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
	// 调试用，若 code 以 "!" 开头，则 OpenID 等于 code 将 "!" 替换为 "+"
	if Config.Debug && len(code) >= 2 && code[0] == '!' {
		return "+" + code[1:]
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
		close(player.Channel)
		// 复制一份
		player = &Player{
			Id:       player.Id,
			Nickname: player.Nickname,
			Avatar:   player.Avatar,
		}
		Players[player.Id] = player
	}
	player.Channel = make(chan interface{})
	player.InRoom = room
	room.People = append(room.People, player)
	DataMutex.Unlock()

	// 建立连接
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	inChannel := make(chan interface{})
	outChannel := player.Channel

	go func(c *websocket.Conn, ch chan interface{}) {
		var object interface{}
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

messageLoop:
	for inChannel != nil && outChannel != nil {
		select {
		case object, ok := <-inChannel:
			if !ok {
				inChannel = nil
				break
			}
			fmt.Println(object)

		case object, ok := <-outChannel:
			if !ok {
				outChannel = nil
				break
			}
			if err := c.WriteJSON(object); err != nil {
				log.Println(err)
				break messageLoop
			}
		}
	}

	c.Close()
	// inChannel 由 goroutine 在 c.Close() 触发的错误之后关闭
	if outChannel != nil {
		close(outChannel)
	}

	DataMutex.Lock()
	player.InRoom.People = removeElement(player.InRoom.People, player)
	player.InRoom = nil
	player.Channel = nil // 清空 Channel 作为连接结束的信号
	DataMutex.Unlock()

	log.Println("connection closed")
}

func SetUpHttp() {
	port := 2310

	http.HandleFunc("/channel/", channelHandler)

	log.Printf("Listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
