package utils

type Message struct {
	Auth *MessageAuth `json:"auth"`
	Data interface{}  `json:"data"`
}
type MessageAuth struct {
	AppID string `json:"app_id"`
	Sign  string `json:"sign"`
}
