package fcm_model

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type FcmEvent struct {
	gorm.Model
	Type  string `json:"type" db:"type" gorm:"type:text"`
	Body  string `json:"body" db:"body" gorm:"type:text"`
	Retry bool   `json:"retry" db:"retry" gorm:"type:text"`
}

func (m *FcmEvent) ToByte() ([]byte, error) {
	data, err := json.Marshal(m)
	return data, err
}
