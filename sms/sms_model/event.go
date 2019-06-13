package sms_model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type SmsEvent struct {
	gorm.Model
	Type  string `json:"type" db:"type" gorm:"type:text"`
	Body  string `json:"body" db:"body" gorm:"type:text"`
	Retry bool   `json:"retry" db:"retry" gorm:"type:text"`
}

func (m *SmsEvent) ToByte() ([]byte, error) {
	data, err := json.Marshal(m)
	return data, err
}
