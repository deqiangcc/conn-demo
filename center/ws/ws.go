package main

import (
	"conn-demo/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
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

	appID, hmac := parseParam(req)
	if hmac == "" || appID == "" {
		log.Println("param err")
		return
	}
	if !utils.HmacVerify(utils.APP_SECRET, hmac) {
		log.Println("invalid hmac: ", hmac)
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