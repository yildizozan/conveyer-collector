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

type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func NewPosition(x, y, z float32) *Position {
	return &Position{X: x, Y: y, Z: z}
}

func (p *Position) MarshallJSON() ([]byte, error) {
	return json.Marshal(&struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
		Z float32 `json:"z"`
	}{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
	})
}

type Engines struct {
	Green  int32 `json:"green"`
	Red    int32 `json:"red"`
	Black  int32 `json:"black"`
	Orange int32 `json:"orange"`
}

func NewEngines(green int32, red int32, black int32, orange int32) *Engines {
	return &Engines{Green: green, Red: red, Black: black, Orange: orange}
}

func (e *Engines) MarshallJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Green  int32 `json:"green"`
		Red    int32 `json:"red"`
		Black  int32 `json:"black"`
		Orange int32 `json:"orange"`
	}{
		Green:  0,
		Red:    0,
		Black:  0,
		Orange: 0,
	})
}
