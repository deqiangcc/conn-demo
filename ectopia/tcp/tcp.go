package main

import (
	"conn-demo/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Fatal("Failed to connect to test server")
	}
	fmt.Println("connect tcp server success ...")

	go read(conn)
	send(conn)
}

func send(conn net.Conn) {
	for {
		var data string
		fmt.Scanln(&data)

		msg := utils.Message{
			Auth: &utils.MessageAuth{
				AppID: APP_ID,
				Sign:  utils.Sign(APP_ID, APP_SECRET),
			},
			Data: data,
		}
		msgJson, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("marshal msg err:", err)
			continue
		}
		_, err = conn.Write(msgJson)
		if err != nil {
			fmt.Println("send msg err:", err)
		}
	}
}

func read(conn net.Conn) {
	for {
		buf := [4096]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Read fail err", err)
			return
		}
		fmt.Println("received msg：", string(buf[:n]))
	}
}
