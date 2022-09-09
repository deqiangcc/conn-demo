package main

import (
	"conn-demo/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	//APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

func main() {
	http.HandleFunc("/ws", wsPage)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Println("start websoket server err: ", err)
		return
	}

	fmt.Println("start websoket server success ...")
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	appID, sgin := parseParam(req)
	if sgin == "" || appID == "" {
		log.Println("param err")
		return
	}
	if !utils.SignVerify(appID, APP_SECRET, sgin) {
		log.Println("invalid sign: ", sgin)
		return
	}

	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	go read(conn)
	go send(conn)
}

func parseParam(req *http.Request) (appID, sgin string) {
	rawQuery := req.URL.RawQuery
	if rawQuery == "" {
		return
	}
	arr := strings.Split(req.URL.RawQuery, "&")
	appID = arr[0]
	sgin = arr[1]

	return
}

func read(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		fmt.Printf("received msg: %s\n", string(message))
	}
}


func send(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		var msg string
		fmt.Scanln(&msg)
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("send msg err:", err)
			return
		}
	}
}