package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	logic.LogFile() // open file app.log

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Хотите запустить TelegramBot? y/any_key\n> ")
	tgbot, _ := reader.ReadString('\n')
	tgbot = strings.TrimSpace(tgbot)
	if tgbot == "y" || tgbot == "н" {
		telegram.Bot(db)
	}

	var menu string
	for {
		output.Hello()
		menu, _ = reader.ReadString('\n')
		menu = strings.TrimSpace(menu)

		switch menu {
		case "1":
			output.WeatherPrint(0)
			var city string
			for {
				city, _ = reader.ReadString('\n')
				city = strings.TrimSpace(city)
				if city == "q" || city == "й" {
					break
				}

				city, temp, conditions, pressure, wind_speed, err := weather.WeatherFunc(city)
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

					notification := database.NotificationConditionsPressureWind_speed(db, temp, conditions, pressure, wind_speed)

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
			}
		case "3":
			var StyleSwitch string
			output.PrintClothingAdviceResult_Hello()
			StyleSwitch, _ = reader.ReadString('\n')
			StyleSwitch = strings.TrimSpace(StyleSwitch)
			StyleSwitchInt, _ := strconv.Atoi(StyleSwitch)
			switch StyleSwitchInt {
			case 1:
				style, StyleString, resstyle, err := database.ClothingAdvice(db, StyleSwitchInt)
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
				}

				style, StyleString, resstyle, err := database.ClothingAdvice(db, StyleSwitchInt)
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
		case "q":
			return
		}
	}
}
