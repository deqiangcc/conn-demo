package main

import (
	"conn-demo/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"time"
)
type request struct {
	Msg string `json:"msg"`
}

var conn *net.TCPConn

func main() {
	var err error
	//dialer := net.Dialer{}
	//dialer.KeepAlive = time.Second * 3
	//conn, err = dialer.Dial("tcp", "127.0.0.1:8001")
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:8001")
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Failed to connect to server")
	}
	fmt.Println("connect tcp server success ...")

	err = conn.SetKeepAlive(true)
	if err != nil {
		fmt.Println("set keepAlive err:", err)
		return
	}
	if err = conn.SetKeepAlivePeriod(3*time.Second); err != nil {
		fmt.Println("set keepAlive time err:", err)
		return
	}

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
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				var receivedMsg utils.CenterMessage
				err = json.Unmarshal(buf[:n], &receivedMsg)
				if err != nil {
					fmt.Println("unmarshal msg err:", err)
					c.AbortWithStatus(http.StatusBadRequest)
					break
				}
				fmt.Printf("received msg: %+v\n", receivedMsg)
				if receivedMsg.Data != nil {
					c.JSON(http.StatusOK, receivedMsg)
					break
				}
			}
		} else {
			fmt.Println("disconnect")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

	})
	router.Run(":8003")
}
