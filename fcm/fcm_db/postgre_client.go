package fcm_db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"os"
)

type PostgresClient struct {
	conn *gorm.DB
}

func (c *PostgresClient) Init() error {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	log.Infoln("connectionStr :  " + connectionStr)
	db, err := gorm.Open("postgres", connectionStr)
	c.conn = db
	return err
}

func (c *PostgresClient) Save(model interface{}) error {
	c.conn.AutoMigrate(model)
	db := c.conn.Save(model)
	return db.Error
}

func (c *PostgresClient) GetAll(list interface{}) error {
	db := c.conn.Find(list)
	return db.Error
}

func (c *PostgresClient) Close() {
	c.conn.Close()
}
