package model

import "encoding/json"

type Event struct {
	Type  string `json:"type"`
	Body  string `json:"body"`
	Retry bool   `json:"retry"`
}

func (m *Event) ToByte() ([]byte, error) {
	data, err := json.Marshal(m)
	return data, err
}
