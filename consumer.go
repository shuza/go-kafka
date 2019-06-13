package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/shuza/go-kafka/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"io/ioutil"
	"net/http"
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
			log.Warnf("Can't parse event  Error :  %v\n", err)
			log.Warnf("Message on %s  :  %s\n", msg.TopicPartition, string(msg.Value))
		}

		switch strings.ToLower(event.Type) {
		case "sms":
			log.Printf("SMS type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
			data, err := event.ToByte()
			if err != nil {
				log.Warnf("can't create sms service request body  Error  :  %v\n", err)
				break
			}
			resp, err := doPostRequest("http://localhost:8081/event", data)
			if err != nil {
				log.Warnf("Failed to connect sms service  Error  :  %v\n", err)
				break
			}
			log.Printf("Sms Service Response  :  %v \t %v \n", resp["status"], resp["message"])
		case "fcm":
			log.Printf("FCM type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
			data, err := event.ToByte()
			if err != nil {
				log.Warnf("can't create fcm service request body  Error  :  %v\n", err)
				break
			}
			resp, err := doPostRequest("http://localhost:8082/event", data)
			if err != nil {
				log.Warnf("Failed to connect fcm service  Error  :  %v\n", err)
				break
			}
			log.Printf("Fcm Service Response  :  %v \t %v \n", resp["status"], resp["message"])
		case "socket":
			log.Printf("Socket type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		default:
			log.Printf("Unknown type notification : %v  %v  %v\n", event.Type, event.Body, event.Retry)
		}
	}

	consumer.Close()
}

func doPostRequest(url string, data []byte) (map[string]interface{}, error) {
	buf := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response := make(map[string]interface{})
	err = json.Unmarshal(body, &response)

	return response, err
}
