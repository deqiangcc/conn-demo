package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {
	var broker = "e1a1d163.cn-shenzhen.emqx.cloud"
	var port = 11427
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("mqtt_client_2")
	opts.SetUsername("qdq")
	opts.SetPassword("qdq123456")
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connect mqtt broker success ...")

	sub(client)
}

func sub(client mqtt.Client) {
	for  {
		qos := 0
		client.Subscribe("mqtt_topic_1", byte(qos), func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		})
	}
}