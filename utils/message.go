package utils

import (
	"encoding/json"
	"fmt"
)

const (
	APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

type Message struct {
	Auth *MessageAuth `json:"auth"`
	Data interface{}  `json:"data"`
}
type MessageAuth struct {
	AppID string `json:"app_id"`
	Sign  string `json:"sign"`
}

func NewMsg(data interface{}) ([]byte, error) {
	msg := Message{
		Auth: &MessageAuth{
			AppID: APP_ID,
			Sign:  Sign(APP_ID, APP_SECRET),
		},
		Data: data,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return nil, err
	}

	return msgJson, err
}