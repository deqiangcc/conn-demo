package main

import (
	"conn-demo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

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
	if err = conn.SetKeepAlivePeriod(10 * time.Second); err != nil {
		fmt.Println("set keepAlive time err:", err)
		return
	}

	msg := &utils.BrokerMessage{
		AppID: utils.APP_ID,
		Type:  utils.BrokerMsgTypeAuthRequest,
		Data:  utils.GenHmac(utils.APP_SECRET),
	}
	err = send(msg)
	if err != nil {
		fmt.Println("send msg err:", err)
		return
	}

	for {
		if err = readBrokerMsg(); err != nil {
			fmt.Println("read msg err:", err)
			return
		}
	}

}

func send(msg *utils.BrokerMessage) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return err
	}
	_, err = conn.Write(msgJson)
	return err
}

func readBrokerMsg() error {
	buf := [4096]byte{}
	n, err := conn.Read(buf[:])
	if err != nil {
		return err
	}
	data := buf[:n]

	var msg utils.BrokerMessage
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	fmt.Printf("received broker msg：%+v\n", msg)
	brokerMsg := &utils.BrokerMessage{
		RequestID: msg.RequestID,
		AppID:     utils.APP_ID,
	}
	switch msg.Type {
	case utils.BrokerMsgTypeAuthResponse:
		fmt.Println("auth ret:", msg.Data)
		return nil
	case utils.BrokerMsgTypeTestRequest:
		brokerMsg.Type = utils.BrokerMsgTypeTestResonse
		if conn != nil {
			brokerMsg.Data = msg.Data
		} else {
			brokerMsg.Data = fmt.Sprintf("invalid msg: %+v", msg)
		}
	default:
		brokerMsg.Type = utils.BrokerMsgTypeError
		brokerMsg.Data = fmt.Sprintf("invalid msg type: %+v", msg.Type)
	}

	if err := send(brokerMsg); err != nil {
		return errors.New(fmt.Sprintf("send msg err:: %+v", msg))
	}

	return nil
}
