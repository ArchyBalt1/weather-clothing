package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		log.Println("Ошибка при .env загрузке", err)
		return
	}

	db, err := database.Init()
	if err != nil {
		log.Println("Ошибка при запуске бд", err)
		return
	}
	defer db.Close()

	logic.LogFile()

	var start string
	fmt.Print("Хотите запустить TelegramBot? y/any_key\n> ")
	fmt.Scan(&start)
	if start == "y" || start == "н" {
		telegram.Bot(db)
	}
	var a string
	output.Hello()
	i := 0
	for {
		output.HelloMenu(&i)
		fmt.Scan(&a)
		switch a {
		case "1":
			var city, conditions string
			var temp, pressure int
			var wind_speed float32
			var cityes string
			output.WeatherPrint(0)
			for {
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					cityes = scanner.Text()
				}
				if cityes == "" {
					continue
				}
				if cityes == "q" || cityes == "й" {
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
					err = database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed)
					if err != nil {
						log.Println("Ошибка при insert запросе", err)
						return
					}

					notification := database.NotificationConditionsPressureWind_speed(db, conditions, pressure, wind_speed)

					signal := output.PrintWeatherResult(city, temp, conditions, notification, wind_speed, pressure)
					if signal == "break" {
						break
					}
					output.WeatherPrint(2)
				}
			}
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
			output.PrintHistoryRecent_requests(FilterSlice)
			for {
				signal := output.PrintHistoryResult(wHistory)
				if signal == "break" {
					break
				} else if signal == "continue" {
					continue
				}
			} // Вывод
		case "3":
			var b int
			output.PrintClothingAdviceResult_Hello()
			fmt.Scan(&b)
			switch b {
			case 1:
				style, StyleString, resstyle, err := database.ClothingAdvice(db, b)
				if err != nil {
					log.Println(err)
					return
				}

				for {
					signal := output.PrintClothingAdviceResult(style, StyleString, resstyle)
					if signal == "break" {
						break
					}
				}
			case 2:
				if err := database.HistoryLimit10(db); err != nil {
					log.Println(err)
					return
				} // фильтруем 10 последних записей

				style, StyleString, resstyle, err := database.ClothingAdvice(db, b)
				if err != nil {
					log.Println(err)
					return
				}

				for {
					if StyleString == nil {
						break
					}
					signal := output.PrintClothingAdviceResult(style, StyleString, resstyle)
					if signal == "break" {
						break
					}
				}
			}
		default:
			return
		}
	}
}
