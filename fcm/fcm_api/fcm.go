package fcm_api

import (
	"github.com/gin-gonic/gin"
	"github.com/shuza/go-kafka/fcm/fcm_db"
	"github.com/shuza/go-kafka/fcm/fcm_model"
	log "github.com/sirupsen/logrus"
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

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": "Fcm event processed successfully",
	})
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
