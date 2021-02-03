package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func channelHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		text := string(msg)
		response := fmt.Sprintf("Text is %s, length is %d", text, len([]rune(text)))
		err = c.WriteMessage(msgType, []byte(response))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func SetUpHttp() {
	port := 2310

	http.HandleFunc("/channel", channelHandler)

	log.Printf("Listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
