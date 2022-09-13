package main

import (
	"conn-demo/utils"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sync"
)

var KeepFlag = []byte{0xFF, 0xFF, 0xFF, 0xFF}

type CenterConnection struct {
	Conn   net.Conn // 连接信息
	IsAuth bool     // 是否已鉴权
}

// 连接集合
//var CenterConnectionMap = make(map[string]*CenterConnection)
var CenterConnectionMap sync.Map

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8001,
	})
	if err != nil {
		fmt.Println("Listen tcp server failed,err:", err)
		return
	}
	fmt.Println("start tcp server success ...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Listen.Accept failed,err:", err)
	}
	for {
		read(conn)
	}
}

func send(conn net.Conn, msg *utils.CenterMessage) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(msgJson)
	return err
}

func read(conn net.Conn) {
	buf := [4096]byte{}
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("Read fail err", err)
		return
	}
	data := buf[:n]

	var msg utils.CenterMessage
	err = json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("unmarshal msg err:", err)
		return
	}

	fmt.Printf("received msg：%+v\n", msg)
	switch msg.Type {
	case utils.CenterMsgTypeAuth:
		if !utils.HmacVerify(utils.APP_SECRET, (msg.Data).(string)) {
			fmt.Println("hmac verify failed")
			return
		}
		CenterConnectionMap.Store(msg.AppID, CenterConnection{
			Conn:   conn,
			IsAuth: true,
		})
	default:
		val, ok := CenterConnectionMap.Load(msg.AppID)
		conn := (val).(CenterConnection)
		if ok {
			if !conn.IsAuth {
				fmt.Println("connect not auth")
				return
			}
			testMsg := &utils.CenterMessage{
				AppID: utils.APP_ID,
				Type: utils.CenterMsgTypeTest,
				Data: "hello world",
			}
			if err := send(conn.Conn, testMsg); err != nil {
				fmt.Println("send msg err:", err)
				return
			}
		} else {
			fmt.Printf("invalid msg: %+v\n", msg)
		}
	}
}

func IsKeep(data []byte) bool {
	if len(data) == 4 {
		return reflect.DeepEqual(data, KeepFlag)
	}
	return false
}
