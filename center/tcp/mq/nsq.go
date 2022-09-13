package mq

import (
	"encoding/json"
	"fmt"
	nsq "github.com/nsqio/go-nsq"
)

const (
	NsqTopicTest = "test"
)

type NsqMessageHandler struct {}

type NsqMessage struct {
	RequestID     string      `json:"request_id"`
	DruidAppID    string      `json:"druid_app_id"`
	BrokerMsgType uint32      `json:"broker_msg_type"`
	Data          interface{} `json:"data"`
}

var NsqProducer *nsq.Producer

func NewNsqProducer() error {
	var err error
	config := nsq.NewConfig()
	NsqProducer, err = nsq.NewProducer("127.0.0.1:4150", config)

	return err
}

func (h *NsqMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	var msg NsqMessage
	if err := json.Unmarshal(m.Body, &msg); err != nil {
		return err
	}
	fmt.Printf("msg:%+v", msg)
	return nil
}

func NsqConsume(topic string) error {
	var err error
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, "channel", config)
	if err != nil {
		return err
	}

	consumer.AddHandler(&NsqMessageHandler{})

	// 建立NSQLookupd连接
	err = consumer.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		return err
	}

	// 建立一个nsqd连接
	if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		return err
	}

	return nil
}

func NsqProduct(topic string, msg *NsqMessage) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err = NsqProducer.Publish(topic, msgByte); err != nil {
		return err
	}
	NsqProducer.Stop()

	return nil
}
