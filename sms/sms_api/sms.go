package sms_api

import (
	"github.com/gin-gonic/gin"
	"github.com/shuza/go-kafka/sms/sms_db"
	"github.com/shuza/go-kafka/sms/sms_model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
)

func AddNewSms(c *gin.Context) {
	var event sms_model.SmsEvent
	if err := c.BindJSON(&event); err != nil {
		log.Warnf("AddNewSms failed to parse event body Error :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})

		return
	}

	if err := sms_db.Client.Init(); err != nil {
		log.Warnf("AddNewSms failed to connect DB  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}
	defer sms_db.Client.Close()

	if err := sms_db.Client.Save(&event); err != nil {
		log.Warnf("AddNewSms failed to save event  ==//  %v", event)
		log.Warnf("AddNewSms failed to save event  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	go sendDeliveryReport(event)

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": "Sms event processed successfully",
	})
}

func sendDeliveryReport(event sms_model.SmsEvent) {
	topic := "deliver-report-sms"
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		log.Warnf("Can't connect to kafka to send deliver report  Error  :  %v\n", err)
		return
	}

	data, _ := event.ToByte()

	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, deliveryChan)
	if err != nil {
		log.Warnf("failed to send SMS deliver report  Error  :  %v\n", err)
		return
	}

	report := <-deliveryChan
	log.Infof("SMS event deliver report send  :  %v\n", report)
}

func AllSmsEvent(c *gin.Context) {
	if err := sms_db.Client.Init(); err != nil {
		log.Warnf("AllSmsEvent failed to connect DB  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}
	defer sms_db.Client.Close()

	events := make([]sms_model.SmsEvent, 0)
	if err := sms_db.Client.GetAll(&events); err != nil {
		log.Warnf("AllSmsEvent failed to get event list  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"data":    events,
		"message": "SUccessful",
	})
}
