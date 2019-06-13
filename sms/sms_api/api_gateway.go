package sms_api

import "github.com/gin-gonic/gin"

func NewGinEngine() *gin.Engine {
	r := gin.Default()

	r.GET("/", Index)
	r.POST("/event", AddNewSms)
	r.GET("/event", AllSmsEvent)

	return r
}
