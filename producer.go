package main

import (
	"bufio"
	"fmt"
	"github.com/shuza/go-kafka/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	go deliveryReport(producer)

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Type your topic : ")
		reader.Scan()
		topic := reader.Text()

		fmt.Print("Type message type : ")
		reader.Scan()
		eventType := reader.Text()

		fmt.Print("Type your message : ")
		reader.Scan()
		message := reader.Text()

		event := model.Event{eventType, message, false}

		if data, err := event.ToByte(); err == nil {
			sendMessage(producer, topic, data)
		}
	}
}

func sendMessage(producer *kafka.Producer, topic string, payload []byte) {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
			Partition:kafka.PartitionAny,
		},
		Value: payload,
	}, nil)

	if err != nil {
		log.Warnf("sendMessage failed : %v\n", string(payload))
	}
}

func deliveryReport(producer *kafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Warnf("Failed to send %v \t Error \t ==:: \t %v \n", ev.TopicPartition, ev.TopicPartition.Error)
			} else {
				log.Printf("Successfully send to %v \n", ev.TopicPartition)
			}
		}
	}
}
