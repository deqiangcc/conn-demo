package main

import (
	"fmt"
	"log"
	"net"
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
		var msg string
		fmt.Scanln(&msg)
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("send msg err:", err)
		}
	}
}

func read(conn net.Conn) {
	for  {
		buf := [4096]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("Read fail err", err)
			return
		}
		fmt.Println("received msg：", string(buf[:n]))
	}
}
