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

var conn net.Conn

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Fatal("Failed to connect to test server")
	}
	fmt.Println("connect tcp server success ...")

	startHttpserver(conn)
	go read(conn)
	send(conn)
}

func send(conn net.Conn) {
	for {
		var data string
		fmt.Scanln(&data)

		msg, err := utils.NewMsg(data)
		if err != nil {
			fmt.Println("new msg err:", err)
			continue
		}
		_, err = conn.Write(msg)
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

type request struct {
	Msg string `json:"msg"`
}

func startHttpserver(conn net.Conn) {
	router := gin.Default()
	router.POST("/test1", func(c *gin.Context) {
		var req request
		if err := c.BindJSON(&req); err != nil {
			fmt.Println("new msg err: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		//fmt.Println(conn)
		//fmt.Println(postMsg)
		if conn != nil {
			sendMsg, err := utils.NewMsg(req.Msg)
			if err != nil {
				fmt.Println("new msg err: ", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			_, err = conn.Write(sendMsg)
			for {
				buf := [4096]byte{}
				n, err := conn.Read(buf[:])
				if err != nil {
					fmt.Println("Read fail err", err)
					return
				}
				var receivedMsg utils.Message
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
