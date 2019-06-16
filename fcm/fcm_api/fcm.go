package fcm_api

import (
	"github.com/gin-gonic/gin"
	"github.com/shuza/go-kafka/fcm/fcm_db"
	"github.com/shuza/go-kafka/fcm/fcm_model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
)

func AddNewFcm(c *gin.Context) {
	var event fcm_model.FcmEvent
	if err := c.BindJSON(&event); err != nil {
		log.Warnf("AddNewFcm failed to parse event body Error :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})

		return
	}

	if err := fcm_db.Client.Init(); err != nil {
		log.Warnf("AddNewFcm failed to connect DB  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}
	defer fcm_db.Client.Close()

	if err := fcm_db.Client.Save(&event); err != nil {
		log.Warnf("AddNewFcm failed to save event  ==//  %v", event)
		log.Warnf("AddNewFcm failed to save event  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}

	go sendDeliveryReport(event)

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": "Fcm event processed successfully",
	})
}

func sendDeliveryReport(event fcm_model.FcmEvent) {
	topic := "deliver-report-fcm"
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
		log.Warnf("failed to send fcm deliver report  Error  :  %v\n", err)
		return
	}

	report := <-deliveryChan
	log.Infof("fcm event deliver report send  :  %v\n", report)
}

func AllFcmEvent(c *gin.Context) {
	if err := fcm_db.Client.Init(); err != nil {
		log.Warnf("AllFcmEvent failed to connect DB  Error  :  %v", err)
		c.JSON(200, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})

		return
	}
	defer fcm_db.Client.Close()

	events := make([]fcm_model.FcmEvent, 0)
	if err := fcm_db.Client.GetAll(&events); err != nil {
		log.Warnf("AllFcmEvent failed to get event list  Error  :  %v", err)
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
