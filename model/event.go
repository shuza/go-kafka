package model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type Event struct {
	Type  string `json:"type"`
	Body  string `json:"body"`
	Retry bool   `json:"retry"`
}

func (m *Event) ToByte() ([]byte, error) {
	data, err := json.Marshal(m)
	return data, err
}

type DeliverReport struct {
	gorm.Model
	Type        string `json:"type" db:"type" gorm:"type:text"`
	Body        string `json:"body" db:"body" gorm:"type:text"`
	Retry       bool   `json:"retry" db:"retry" gorm:"type:text"`
	IsDelivered bool   `json:"is_delivered" db:"is_delivered" gorm:"type:text"`
}
