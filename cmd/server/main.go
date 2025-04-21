package main

import (
	"fmt"
	"log"
	database "weather-clothing/internal/db"
	"weather-clothing/internal/output"
	w "weather-clothing/internal/weather"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../../.env") // Достали ключ
	if err != nil {
		log.Println("Ошибка при .env загрузке")
		return
	}

	db, err := database.Init()
	if err != nil {
		log.Println("Ошибка при запуске бд")
		return
	}
	defer db.Close()

	var a string
	for {
		output.Hello()
		fmt.Scan(&a)
		switch a {
		case "1":
			city, temp, conditions, pressure, wind_speed, err := w.WeatherFunc()
			if city == "Arpol" {
				continue
			} else if err != nil {
				log.Println("Ошибка при получении погодных условий", err)
				return
			}

			err = database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed)
			if err != nil {
				log.Println("Ошибка при insert запросе", err)
				return
			}

			notification := database.NotificationConditionsPressureWind_speed(db, conditions, pressure, wind_speed)

			output.PrintWeatherResult(city, temp, conditions, notification, wind_speed, pressure)
		case "2":
			err := database.ReadHistory(db) // вывод в функции
			if err != nil {
				log.Println(err)
				return
			}
		case "3":
			err := database.ClothingAdvice(db)
			if err != nil {
				log.Println(err)
				return
			}
		default:
			return
		}
	}
}
