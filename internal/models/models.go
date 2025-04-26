package models

import "time"

type WeatherHistory_10 struct {
	City       string
	Temp       int
	Date       time.Time
	Conditions string
	Pressure   int
	Wind_speed float32
} // Хранение 10 последних записей

type ResStyle struct {
	Style       string
	Comments    string
	Accessories string
} // Хранится итоговый выбранный стиль

type CityStyle struct {
	City       string
	Temp       int
	Conditions string
	Wind_speed float32
} // Хранится город, под данные которого будет подобран стиль
