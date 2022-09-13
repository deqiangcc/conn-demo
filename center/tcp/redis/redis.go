package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

type CenterMessage struct {
	RequestID     string      `json:"request_id"`
	DruidAppID    string      `json:"druid_app_id"`
	BrokerMsgType uint32      `json:"broker_msg_type"`
	Data          interface{} `json:"data"`
}


func ConnRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func SetCenterRequestMsg(msg *CenterMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = redisClient.HSetNX("center-request-messages", msg.RequestID, string(msgBytes)).Result()

	return err
}

func DelCenterRequestMsg(msg *CenterMessage) error {
	_, err := redisClient.HDel("center-request-messages", msg.RequestID).Result()

	return err
}

func SetCenterResponseMsg(msg *CenterMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = redisClient.HSetNX("center-response-messages", msg.RequestID, string(msgBytes)).Result()

	return err
}

func DelCenterReponseMsg(msg *CenterMessage) error {
	_, err := redisClient.HDel("center-response-messages", msg.RequestID).Result()

	return err
}

func GetCenterRequestMsg(requestID string) (*CenterMessage, error) {
	ret, err := redisClient.HGet("center-request-messages", requestID).Result()
	if err != nil {
		return nil, err
	}
	var msg *CenterMessage
	if err = json.Unmarshal([]byte(ret), &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func GetCenterResponseMsg(requestID string) (*CenterMessage, error) {
	ret, err := redisClient.HGet("center-response-messages", requestID).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	var msg *CenterMessage
	if err = json.Unmarshal([]byte(ret), &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func GetCenterRequestMsgAll() ([]*CenterMessage, error) {
	ret, err := redisClient.HGetAll("center-request-messages").Result()
	if err != nil {
		return nil, err
	}
	msgs := []*CenterMessage{}
	for _, val := range ret {
		var msg *CenterMessage
		if err = json.Unmarshal([]byte(val), &msg); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}


	return msgs, nil
}