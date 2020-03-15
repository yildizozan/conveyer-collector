package model

import "encoding/json"

type Measurement struct {
	Weight   float32 `json:"weight"`
	Humidity float32 `json:"humidity"`
	Color    float32 `json:"color"`
}

func NewMeasurement(weight float32, humidity float32, color float32) *Measurement {
	return &Measurement{Weight: weight, Humidity: humidity, Color: color}
}

func (m *Measurement) MarshallJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Weight   float32 `json:"weight"`
		Humidity float32 `json:"humidity"`
		Color    float32 `json:"color"`
	}{
		Weight:   m.Weight,
		Humidity: m.Humidity,
		Color:    m.Color,
	})
}
