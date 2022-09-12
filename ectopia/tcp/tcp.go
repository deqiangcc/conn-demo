package main

import (
	"conn-demo/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
)
type request struct {
	Msg string `json:"msg"`
}

var conn net.Conn

func main() {
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Fatal("Failed to connect to test server")
	}
	fmt.Println("connect tcp server success ...")

	msg := &utils.CenterMessage{
		AppID: utils.APP_ID,
		Type: utils.CenterMsgTypeAuth,
		Data: utils.GenHmac(utils.APP_SECRET),
	}
	err = send(msg)
	if err != nil {
		fmt.Println("send msg err:", err)
		return
	}

	startHttpserver()
	go read()
}


func send(msg *utils.CenterMessage) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return err
	}
	_, err = conn.Write(msgJson)
	return err
}

func read() {
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

func startHttpserver() {
	router := gin.Default()
	router.POST("/test1", func(c *gin.Context) {
		var req request
		if err := c.BindJSON(&req); err != nil {
			fmt.Println("new msg err: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if conn != nil {
			sendMsg := &utils.CenterMessage{
				AppID: utils.APP_ID,
				Type: utils.CenterMsgTypeTest,
				Data: req.Msg,
			}
			if err := send(sendMsg); err != nil {
				fmt.Println("send msg err", err)
				return
			}
			for {
				buf := [4096]byte{}
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Println("Read fail err", err)
					return
				}
				var receivedMsg utils.CenterMessage
				err = json.Unmarshal([]byte(string(buf[:n])), &receivedMsg)
				if err != nil {
					fmt.Println("unmarshal msg err:", err)
					break
				}
				fmt.Printf("received msg: %+v\n", receivedMsg)
				if receivedMsg.Data != nil {
					c.JSON(http.StatusOK, receivedMsg)
					break
				}
			}
		} else {
			c.String(http.StatusOK, "can not find gwid,please check gwid")
		}

	})
	router.Run(":8083")
}
