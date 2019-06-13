package fcm_db

type IFcmDb interface {
	Init() error
	Save(model interface{}) error
	GetAll(list interface{}) error
	Close()
}

var Client IFcmDb
