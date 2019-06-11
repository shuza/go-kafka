package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{"test", "^aRegex.*[Tt]opic"}, nil)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Warnf("Consumer Error \t ==:: \t %v  (%v)\n", err, msg)
			continue
		}

		log.Printf("Message on %s  :  %s\n", msg.TopicPartition, string(msg.Value))
	}

	consumer.Close()
}
