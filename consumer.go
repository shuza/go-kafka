package main

import (
	"encoding/json"
	"flag"
	"github.com/shuza/go-kafka/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "bootstrap server list")
	consumerGroup := flag.String("group", "myGroup", "consumer group name")
	topic := flag.String("topic", "test", "topic you want to consume")
	flag.Parse()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": *host,
		"group.id":          *consumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{*topic, "^aRegex.*[Tt]opic"}, nil)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Warnf("Consumer Error \t ==:: \t %v  (%v)\n", err, msg)
			continue
		}
		var event model.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Warnf("Can't parse event  Error :  %v", err)
			log.Warnf("Message on %s  :  %s\n", msg.TopicPartition, string(msg.Value))
		}

		switch strings.ToLower(event.Type) {
		case "sms":
			log.Printf("SMS type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		case "fcm":
			log.Printf("FCM type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		case "socket":
			log.Printf("Socket type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		default:
			log.Printf("Unknown type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		}
	}

	consumer.Close()
}
