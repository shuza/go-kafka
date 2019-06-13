package fcm_api

import "github.com/gin-gonic/gin"

func NewGinEngine() *gin.Engine {
	r := gin.Default()

	r.GET("/", Index)
	r.POST("/event", AddNewFcm)
	r.GET("/event", AllFcmEvent)

	return r
}
