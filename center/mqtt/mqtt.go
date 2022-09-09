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
	opts.SetClientID("mqtt_client_1")
	opts.SetUsername("qdq")
	opts.SetPassword("qdq123456")
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connect mqtt broker success ...")

	pub(client)
}

func pub(client mqtt.Client) {
	for  {
		var msg string
		fmt.Scanln(&msg)
		if msg != "" {
			token := client.Publish("mqtt_topic_1", 0, false, msg)
			token.Wait()
		}
	}
}