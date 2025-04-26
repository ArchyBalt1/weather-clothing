package models

import "time"

type WeatherHistory_10 struct {
	City       string
	Temp       int
	Date       time.Time
	Conditions string
	Pressure   int
	Wind_speed float32
}

type ResStyle struct {
	Style       string
	Comments    string
	Accessories string
}

type Style struct {
	City       string
	Temp       int
	Conditions string
	Wind_speed float32
}
