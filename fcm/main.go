package main

import (
	"github.com/shuza/go-kafka/fcm/fcm_api"
	"github.com/shuza/go-kafka/fcm/fcm_db"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	initFcmDb()
	r := fcm_api.NewGinEngine()

	port := os.Getenv("APP_PORT")
	log.Printf("Fcm service is running on port  %v\n", port)

	r.Run(port)
}

func initFcmDb() {
	fcm_db.Client = &fcm_db.PostgresClient{}
}
