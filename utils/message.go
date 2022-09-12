package utils

const (
	APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

const (
	CenterMsgTypeAuth = 1 // 鉴权
	CenterMsgTypeTest = 2 // 测试
)

type CenterMessage struct {
	AppID string      `json:"app_id"` // 平台id
	Type  uint32      `json:"type"`   // 消息类型：1-auth：受权请求，2-gps request，3...
	Data  interface{} `json:"data"`   // 数据
}
