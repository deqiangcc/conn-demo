package main

import (
	"conn-demo/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

const (
	APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

func main() {
	u := url.URL{
		Scheme:   "ws",
		Host:     "127.0.0.1:8001",
		Path:     "/ws",
		RawQuery: fmt.Sprintf("%s&%s", APP_ID, utils.GenToken(APP_ID, APP_SECRET)),
	}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read msg err:", err)
			return
		}

		fmt.Printf("received msg: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		var name string
		fmt.Scanln(&name)
		err := conn.WriteMessage(websocket.TextMessage, []byte(name))
		if err != nil {
			fmt.Println("send msg err:", err)
			return
		}
	}
}
