package main

import (
	"conn-demo/utils"
	"encoding/json"
	"errors"
	"fmt"
	ts "github.com/0987363/tcp_server"
	"sync"
	"time"
)

type CenterConnection struct {
	Conn   *ts.Context
	IsAuth bool // 是否已鉴权
}

// 连接集合
//var CenterConnectionMap = make(map[string]*CenterConnection)
var CenterConnectionMap sync.Map

func main() {
	server := ts.New("127.0.0.1:8001")
	server.SetUdpProc(10)
	server.SetTimeout(time.Second * 30)
	server.SetCacheSize(4096)
	server.OnNewMessage(func(c *ts.Context) {
		read(c)
		go sendMsg(c)
	})

	fmt.Println("start tcp server success ...")
	server.Listen()
}

func sendMsg(c *ts.Context) {
	for {
		var data string
		fmt.Scanln(&data)
		msg := &utils.CenterMessage{
			Type: utils.CenterMsgTypeTest,
			Data: data,
		}
		err := send(c, msg)
		if err != nil {
			fmt.Println("send msg err:", err)
		}
	}
}

func send(c *ts.Context, msg *utils.CenterMessage) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.Send(msgJson)

	return err
}

func read(c *ts.Context) {
	data := c.ReadData()
	defer c.Trim(len(data))

	var msg utils.CenterMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("unmarshal msg err:", err)
		c.AbortWithError(err)
		return
	}

	fmt.Printf("received msg：%+v\n", msg)
	switch msg.Type {
	case utils.CenterMsgTypeAuth:
		if !utils.HmacVerify(utils.APP_SECRET, (msg.Data).(string)) {
			fmt.Println("hmac verify failed")
			c.AbortWithError(errors.New("hmac verify err"))
			return
		}
		CenterConnectionMap.Store(msg.AppID, CenterConnection{
			Conn: c,
			IsAuth: true,
		})
	default:
		val, ok := CenterConnectionMap.Load(msg.AppID)
		conn := (val).(CenterConnection)
		if ok {
			if !conn.IsAuth {
				fmt.Println("connect not auth")
				c.AbortWithError(errors.New("connect not auth"))
				return
			}
			testMsg := &utils.CenterMessage{
				Type: utils.CenterMsgTypeTest,
				Data: "hello world",
			}
			if err := send(conn.Conn, testMsg); err != nil {
				fmt.Println("send msg err:", err)
				c.AbortWithError(err)
				return
			}
		} else {
			fmt.Printf("invalid msg: %+v\n", msg)
		}
	}
}
