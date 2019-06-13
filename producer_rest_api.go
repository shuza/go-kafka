package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shuza/go-kafka/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Producer rest api is Running...",
		})
	})

	r.POST("/kafka/:topic", createNewMessage)
	r.Run(":8080")
}

func createNewMessage(c *gin.Context) {
	var event model.Event
	if err := c.BindJSON(&event); err != nil {
		log.Warnf("can't parse request body  Error  :  %v\n", err)
		c.JSON(200, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Can't parse request body",
			"data":    err.Error(),
		})

		return
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		log.Warnf("Can't connect to kafka  Error  :  %v\n", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "can't connect to kafka",
			"data":    err.Error(),
		})
	}

	data, _ := event.ToByte()
	deliveryChan := make(chan kafka.Event)
	sendMessageForDelivery(producer, c.Param("topic"), data, deliveryChan)

	report := <-deliveryChan

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": "Successful",
		"data":    report.String(),
	})
}

func sendMessageForDelivery(producer *kafka.Producer, topic string, payload []byte, deliveryChan chan kafka.Event) {
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, deliveryChan)

	if err != nil {
		log.Warnf("sendMessage failed : %v\n", string(payload))
	}
}
