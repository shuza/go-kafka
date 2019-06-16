package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/shuza/go-kafka/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
)

func main() {
	db, err := getDbConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go consumeDeliverReport("deliver-report-fcm", db)
	go consumeDeliverReport("deliver-report-sms", db)

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

func consumeDeliverReport(topic string, db *gorm.DB) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "fcmConsumer",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Warnf("Consumer at %v Error \t ==:: \t %v  (%v)\n", topic, err, msg)
			continue
		}

		var event model.DeliverReport
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Warnf("Can't parse event  Error :  %v\n", err)
			log.Warnf("Message on %s  :  %s\n", msg.TopicPartition, string(msg.Value))
			continue
		}
		event.IsDelivered = true

		db.AutoMigrate(event)
		if err := db.Save(&event).Error; err != nil {
			log.Warnf("Topic : %s  event report save Error  :  %v\n", topic, err)
		}
	}
}

func getDbConnection() (*gorm.DB, error) {
	connectionStr := "postgres://admin:123456@localhost:5432/kafka_db?sslmode=disable"
	db, err := gorm.Open("postgres", connectionStr)
	return db, err
}
