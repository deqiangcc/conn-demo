package main

import (
	"conn-demo/center/tcp/redis"
	"conn-demo/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type request struct {
	Data string `json:"data"`
}

func main() {
	if err := redis.ConnRedis(); err != nil {
		log.Fatal("conn redis err:", err)
	}
	router := gin.Default()
	router.POST("/test1", test)
	router.Run(":8003")
}

func test(c *gin.Context) {
	var req request
	if err := c.BindJSON(&req); err != nil {
		log.Println("param err: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	msg := &redis.CenterMessage{
		RequestID:     utils.RandString(12),
		DruidAppID:    utils.APP_ID,
		BrokerMsgType: utils.BrokerMsgTypeTestRequest,
		Data:          req.Data,
	}
	err := redis.SetCenterRequestMsg(msg)
	if err != nil {
		log.Println("set redis err: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for  {
		responseMsg, err := redis.GetCenterResponseMsg(msg.RequestID)
		if err != nil {
			log.Println("get redis err: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if responseMsg == nil {
			continue
		}

		if err = redis.DelCenterReponseMsg(msg); err != nil {
			log.Println("del redis err: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, responseMsg)
		return
	}

}