package main

import (
	"fmt"
	"log"
	database "weather-clothing/internal/db"
	w "weather-clothing/internal/weather"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("F:\\got\\weather-clothing\\.env") // Достали ключ
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
		fmt.Println("Выберите действие\n1: Узнать температуру\n2: Прочитать бд (последние 10 записей)")
		fmt.Scan(&a)
		switch a {
		case "1":
			city, temp, err := w.WeatherFunc()
			if err != nil {
				log.Println(err)
			}
			fmt.Println(city, temp)
			err = database.WeatherHistory(db, city, temp)
			if err != nil {
				log.Println("Ошибка при insert запросе", err)
				return
			}
		case "2":
			city := "qwe"
			temp := 12
			err := database.ReadHistory(db, city, temp)
			if err != nil {
				log.Println(err)
			}
		default:
			return
		}
	}
}
