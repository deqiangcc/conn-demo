package main

import (
	"conn-demo/center/tcp/redis"
	"conn-demo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type BrokerConnection struct {
	Conn   net.Conn // 连接信息
	IsAuth bool     // 是否已鉴权
}

// 连接集合
//var BrokerConnectionMap = make(map[string]*BrokerConnection)
var BrokerConnectionMap sync.Map

func main() {
	if err := redis.ConnRedis(); err != nil {
		log.Fatal("conn redis err:", err)
	}

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
		fmt.Println("listen accept err:", err)
	}
	go ReadCenterMsg()
	ListenBrokerMsg(conn)
}

func ListenBrokerMsg(conn net.Conn)  {
	for {
		msg, err := readBrokerMsg(conn)
		if err != nil {
			fmt.Println("read broker msg err:", err)
			break
		}
		if err := msgHandle(msg, conn); err != nil {
			fmt.Println("msg handle err:", err)
			break
		}
	}
}

func ReadCenterMsg() {
	for {
		msgs, err := redis.GetCenterRequestMsgAll()
		if err != nil {
			log.Fatal("get center request msg all err:", err)
		}
		if len(msgs) == 0 {
			continue
		}
		for _, msg := range msgs {
			fmt.Printf("received center msg：%+v\n", msg)
			val, ok := BrokerConnectionMap.Load(msg.AppID)
			if !ok {
				log.Println("app disconnect, app_id:", msg.AppID)
				continue
			}
			cconn := (val).(BrokerConnection)
			err = sendDruidPlatformMsg(cconn.Conn, msg)
			if err != nil {
				log.Println("send app msg err:", err)
				continue
			}
			if err = redis.DelCenterRequestMsg(msg); err != nil {
				log.Println("del center request msg err:", err)
				continue
			}
		}

		time.Sleep(time.Second)
	}
}

// 开启http服务
func httpServer() {
	router := gin.Default()
	router.POST("/msg", brokerMsg)
	router.Run(":8002")
}

// 处理web端发送的broker消息
func brokerMsg(ctx *gin.Context) {
	var req utils.BrokerMessageRequest
	if err := ctx.BindJSON(&req); err != nil {
		log.Println("param err:", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	val, ok := BrokerConnectionMap.Load(req.ThirdPlatformAppID)
	if !ok {
		log.Println("app disconnect")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	cconn := (val).(BrokerConnection)
	if !cconn.IsAuth {
		log.Println("app connect not auth")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	msg := &utils.BrokerMessage{
		AppID: req.ThirdPlatformAppID,
		Type:  req.Type,
		Data:  req.Data,
	}
	fmt.Printf("received center msg：%+v\n", msg)
	if err := sendDruidPlatformMsg(cconn.Conn, msg); err != nil {
		log.Println("send  connect not auth")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// todo::超时控制
	for {
		appMsg, err := readBrokerMsg(cconn.Conn)
		if err != nil {
			log.Println("unmarshal msg err:", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			break
		}
		fmt.Printf("received msg: %+v\n", appMsg)
		if appMsg.Data != nil {
			ctx.JSON(http.StatusOK, appMsg)
			break
		}
	}
}

func sendDruidPlatformMsg(conn net.Conn, msg *utils.BrokerMessage) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(msgJson)
	return err
}

func readBrokerMsg(conn net.Conn) (*utils.BrokerMessage, error) {
	buf := [4096]byte{}
	n, err := conn.Read(buf[:])
	if err != nil {
		return nil, err
	}
	data := buf[:n]

	var msg utils.BrokerMessage
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func msgHandle(msg *utils.BrokerMessage, conn net.Conn) error {
	sendMsg := &utils.BrokerMessage{}
	switch msg.Type {
	case utils.BrokerMsgTypeAuthRequest:
		if !utils.HmacVerify(utils.APP_SECRET, (msg.Data).(string)) {
			fmt.Println("hmac verify failed")
			return errors.New("hmac verify failed")
		}
		BrokerConnectionMap.Store(msg.AppID, BrokerConnection{
			Conn:   conn,
			IsAuth: true,
		})
		sendMsg.Type = utils.BrokerMsgTypeAuthResponse
		sendMsg.Data = "auth success"
	case utils.BrokerMsgTypeTestResonse:
		//sendMsg.Type = utils.BrokerMsgTypeTestResonse
		//sendMsg.Data = msg.Data
		if err := redis.SetCenterResponseMsg(msg); err != nil {
			fmt.Println("set center response msg err:", err)
			return err
		}

		return nil

	default:
		sendMsg = &utils.BrokerMessage{
			AppID: utils.APP_ID,
			Type:  utils.BrokerMsgTypeError,
			Data:  "invalid msg type",
		}
	}

	val, ok := BrokerConnectionMap.Load(msg.AppID)
	cconn := (val).(BrokerConnection)
	if ok {
		if !cconn.IsAuth {
			return errors.New("connect not auth")
		}
		if err := sendDruidPlatformMsg(cconn.Conn, sendMsg); err != nil {
			return errors.New(fmt.Sprint("isend msg err:", msg))
		}
	} else {
		return errors.New(fmt.Sprintf("invalid msg: %+v", msg))
	}

	return nil
}
