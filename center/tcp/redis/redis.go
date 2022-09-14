package redis

import (
	"conn-demo/utils"
	"encoding/json"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

//type CenterMessage struct {
//	RequestID     string      `json:"request_id"`
//	DruidAppID    string      `json:"druid_app_id"`
//	BrokerMsgType uint32      `json:"broker_msg_type"`
//	Data          interface{} `json:"data"`
//}

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

const (
	BrokerToCenterMsg = "center_to_broker_msg"
	CenterToBrokerMsg = "broker_to_center_msg"
)

func SetCenterRequestMsg(msg *utils.BrokerMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = redisClient.HSetNX(CenterToBrokerMsg, msg.RequestID, string(msgBytes)).Result()

	return err
}

func DelCenterRequestMsg(msg *utils.BrokerMessage) error {
	_, err := redisClient.HDel(CenterToBrokerMsg, msg.RequestID).Result()

	return err
}

func SetCenterResponseMsg(msg *utils.BrokerMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = redisClient.HSetNX(BrokerToCenterMsg, msg.RequestID, string(msgBytes)).Result()

	return err
}

func DelCenterReponseMsg(msg *utils.BrokerMessage) error {
	_, err := redisClient.HDel(BrokerToCenterMsg, msg.RequestID).Result()

	return err
}

func GetCenterRequestMsg(requestID string) (*utils.BrokerMessage, error) {
	ret, err := redisClient.HGet(CenterToBrokerMsg, requestID).Result()
	if err != nil {
		return nil, err
	}
	var msg *utils.BrokerMessage
	if err = json.Unmarshal([]byte(ret), &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func GetCenterResponseMsg(requestID string) (*utils.BrokerMessage, error) {
	ret, err := redisClient.HGet(BrokerToCenterMsg, requestID).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	var msg *utils.BrokerMessage
	if err = json.Unmarshal([]byte(ret), &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func GetCenterRequestMsgAll() ([]*utils.BrokerMessage, error) {
	ret, err := redisClient.HGetAll(CenterToBrokerMsg).Result()
	if err != nil {
		return nil, err
	}
	msgs := []*utils.BrokerMessage{}
	for _, val := range ret {
		var msg *utils.BrokerMessage
		if err = json.Unmarshal([]byte(val), &msg); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
