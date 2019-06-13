package sms_db

type ISmsDb interface {
	Init() error
	Save(model interface{}) error
	GetAll(list interface{}) error
	Close()
}

var Client ISmsDb
