package main

import (
	"fmt"
	"log"
	database "weather-clothing/internal/db"
	"weather-clothing/internal/logic"
	"weather-clothing/internal/output"
	"weather-clothing/internal/telegram"
	"weather-clothing/internal/weather"

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

	telegram.Bot(db)
	var a string
	for {
		output.Hello()
		fmt.Scan(&a)
		switch a {
		case "1":
			var city, conditions string
			var temp, pressure int
			var wind_speed float32
			var cityes string
			for {
				output.WeatherPrint(0)
				fmt.Scan(&cityes)
				if cityes == "q" {
					break
				}
				city, temp, conditions, pressure, wind_speed, err = weather.WeatherFunc(cityes)
				if city == "Введён неккоректный город" {
					output.WeatherPrint(1)
					continue
				} else if err != nil {
					log.Println("Ошибка при получении погодных условий", err)
					return
				} else {
					break
				}
			}

			err = database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed)
			if err != nil {
				log.Println("Ошибка при insert запросе", err)
				return
			}

			notification := database.NotificationConditionsPressureWind_speed(db, conditions, pressure, wind_speed)

			output.PrintWeatherResult(city, temp, conditions, notification, wind_speed, pressure)
		case "2":
			if err := database.HistoryLimit10(db); err != nil {
				log.Println(err)
				return
			} // фильтруем 10 последних записей
			Slicecity, wHistory, err := database.ReadHistory(db)
			if err != nil {
				log.Println(err)
				return
			} // логика выборки
			FilterSlice := logic.FilterMap(Slicecity, wHistory)
			for {
				signal := output.PrintHistoryResult(FilterSlice, wHistory)
				if signal == "break" {
					break
				}
			} // Вывод
		case "3":
			if err := database.HistoryLimit10(db); err != nil {
				log.Println(err)
				return
			} // фильтруем 10 последних записей
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
