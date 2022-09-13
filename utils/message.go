package utils

const (
	APP_ID     = "6319ac9705b6059fa59de161"
	APP_SECRET = "6319ac97-52fd-fc07-2182-654f-163f5f0f"
)

const (
	BrokerMsgTypeError        = 0 // 错误消息
	BrokerMsgTypeAuthRequest  = 1 // 鉴权请求
	BrokerMsgTypeAuthResponse = 2 // 鉴权响应
	BrokerMsgTypeTestRequest  = 3 // 测试请求
	BrokerMsgTypeTestResonse  = 4 // 测试响应
)

type BrokerMessage struct {
	RequestID string      `json:"request_id"` // 请求id
	AppID     string      `json:"app_id"`     // 平台id
	Type      uint32      `json:"type"`       // 消息类型：1-auth：鉴权请求，2-鉴权响应，3-测试请求，4-测试响应
	Data      interface{} `json:"data"`       // 数据
}

type BrokerMessageRequest struct {
	ThirdPlatformAppID string      `json:"third_platform_app_id"`
	Type               uint32      `json:"type"`
	Data               interface{} `json:"data"`
}

type BrokerMessageResponse struct {
	Data interface{} `json:"data"`
}
