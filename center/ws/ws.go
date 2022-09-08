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
	fmt.Println("start websoket server...")
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":8001", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	appID, token := parseParam(req)
	if token == "" || appID == "" {
		log.Println("param err")
		return
	}
	if !utils.TokenVerify(appID, APP_SECRET, token) {
		log.Println("invalid token: ", token)
		return
	}

	conn, error := upgrader.Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}

	go read(conn)
	go write(conn)
}

func parseParam(req *http.Request) (appID, token string) {
	rawQuery := req.URL.RawQuery
	if rawQuery == "" {
		return
	}
	arr := strings.Split(req.URL.RawQuery, "&")
	appID = arr[0]
	token = arr[1]

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


func write(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

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