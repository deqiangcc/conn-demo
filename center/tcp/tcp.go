package main

import (
	"conn-demo/utils"
	"encoding/json"
	"errors"
	"fmt"
	ts "github.com/0987363/tcp_server"
	"time"
)

const (
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

func main() {
	server := ts.New("127.0.0.1:8001")
	server.SetUdpProc(10)
	server.SetTimeout(time.Second * 30)
	server.SetCacheSize(4096)
	server.OnNewMessage(func(c *ts.Context) {
		read(c)
		go send(c)
	})
	fmt.Println("start tcp server success ...")
	server.Listen()
}

func send(c *ts.Context) {
	for {
		var msg string
		fmt.Scanln(&msg)
		_, err := c.Send([]byte(msg))
		if err != nil {
			fmt.Println("send msg err:", err)
		}
	}
}

func read(c *ts.Context) {
	data := c.ReadData()
	defer c.Trim(len(data))

	var msg utils.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("unmarshal msg err:", err)
		c.AbortWithError(err)
		return
	}
	if err := Verify(&msg); err != nil {
		fmt.Println("msg auth need")
		c.AbortWithError(errors.New("msg auth need"))
		return
	}

	fmt.Printf("received msg：%+v", msg)
}

func Verify(msg *utils.Message) error {
	if msg.Auth == nil {
		fmt.Println("msg auth need")
		return errors.New("msg auth need")
	}

	if !utils.SignVerify(msg.Auth.AppID, APP_SECRET, msg.Auth.Sign) {
		fmt.Println("invalid sign: ", msg.Auth.Sign)
		return errors.New(fmt.Sprint("invalid sign: ", msg.Auth.Sign))
	}

	return nil
}
