package mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateClient(clientID string, handler mqtt.MessageHandler) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(Broker)
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(handler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	} 
	fmt.Println("Connected to MQTT broker")
	return client
}

func Subscribe(client mqtt.Client, topic string, qos byte)  {
	if token := client.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
}