package main

import (
	"github.com/shuza/go-kafka/sms/sms_api"
	"github.com/shuza/go-kafka/sms/sms_db"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	initSmsDb()
	r := sms_api.NewGinEngine()

	port := os.Getenv("APP_PORT")
	log.Printf("Sms service is running on port  %v\n", port)

	r.Run(port)
}

func initSmsDb() {
	sms_db.Client = &sms_db.PostgresClient{}
}
