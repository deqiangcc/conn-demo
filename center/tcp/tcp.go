package main

import (
	"fmt"
	ts "github.com/0987363/tcp_server"
	"time"
)

func main()  {
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

func send(c *ts.Context)  {
	for  {
		var msg string
		fmt.Scanln(&msg)
		_, err := c.Send([]byte(msg))
		if err != nil {
			fmt.Println("send msg err:", err)
		}
	}
}

func read(c *ts.Context)  {
	data := c.ReadData()
	defer c.Trim(len(data))
	fmt.Println("received msg：", string(data))
}
