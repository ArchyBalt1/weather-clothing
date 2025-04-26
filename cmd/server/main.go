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
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Ошибка при загрузке .env файла: %v", err)
		return
	} // Достаём всё необходимое из .env файла

	db, err := database.Init()
	if err != nil {
		log.Printf("Ошибка при подключении к базе данных: %v", err)
		return
	} // Подключение к бд
	defer db.Close()

	logic.LogFile() // открытие файла с логами (app.log)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Хотите запустить TelegramBot? y/any_key\n> ")
	tgbot, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка чтения ввода: %v", err)
		return
	}
	tgbot = strings.TrimSpace(tgbot)

	if tgbot == "y" || tgbot == "н" {
		telegram.Bot(db)
	}

	var menu string
	for {
		output.Hello() // Вывод меню
		menu, _ = reader.ReadString('\n')
		menu = strings.TrimSpace(menu)

		switch menu {
		case "1":
			output.WeatherPrint(0) // В зависимости от цифры разное сообщение
			var city string
			for {
				city, _ = reader.ReadString('\n')
				city = strings.TrimSpace(city)
				if city == "q" || city == "й" {
					break
				}

				city, temp, conditions, pressure, wind_speed, err := weather.WeatherFunc(city) // Сведения о погоде
				if city == "Введён неккоректный город" {
					output.WeatherPrint(1)
					continue
				} else if err != nil {
					log.Printf("Ошибка при получении погодных условий: %v", err)
					return
				} else {
					err = database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed) // Запись в бд
					if err != nil {
						log.Printf("Ошибка при insert запросе: %v", err)
						return
					}

					notification := database.NotificationConditionsPressureWind_speed(db, temp, conditions, pressure, wind_speed) // Достаём из бд советы относительно погодных условий

					signal := output.PrintWeatherResult(city, temp, conditions, notification, wind_speed, pressure) // Печатаем результат
					if signal == "break" {
						break
					}
					output.WeatherPrint(2)
				}
			}
		case "2":
			if err := database.HistoryLimit10(db); err != nil {
				log.Printf("Ошибка при фильтрации 10 элементов: %v", err)
				return
			} // Фильтруем 10 последних записей в бд

			Slicecity, wHistory, err := database.ReadHistory(db)
			if err != nil {
				log.Printf("Ошибка при чтении данных для отображения истории бд: %v", err)
				return
			} // Чтении данных для отображения истории бд

			FilterSlice := logic.FilterMap(Slicecity, wHistory) // Фильтрация для вывода
			output.PrintHistoryRecent_requests(FilterSlice)     // Выводим список бд

			for {
				signal := output.PrintHistoryResult(wHistory) // Печатаем итоговый результат
				if signal == "break" {
					break
				} else if signal == "continue" {
					continue
				} // Различные сигналы для грамотной работы
			}
		case "3":
			var StyleSwitch string
			var StyleSwitchInt int
			output.PrintClothingAdviceResult_Hello()
			for {
				StyleSwitch, err = reader.ReadString('\n')
				if err != nil {
					log.Fatalf("Ошибка чтения ввода: %v", err)
					return
				}
				StyleSwitch = strings.TrimSpace(StyleSwitch)

				StyleSwitchInt, err = strconv.Atoi(StyleSwitch)
				if err != nil {
					fmt.Print("> ")
					continue
				}
				break
			}
			switch StyleSwitchInt { // В зависимости от того, берём мы последний город или выбираем из бд
			case 1:
				style, StyleString, resstyle, err := database.ClothingAdvice(db, StyleSwitchInt) // В данном значении case: Выборка последнего города и стиля под него
				if err != nil {
					log.Printf("Ошибка при выборе города и стиля (case1): %v", err)
					return
				}

				for {
					signal := output.PrintClothingAdviceResult(style, StyleString, resstyle) // Печатаем итоговый результат
					if signal == "break" {
						break
					}
				}
			case 2: // Выбираем из списка
				if err := database.HistoryLimit10(db); err != nil {
					log.Printf("Ошибка при фильтрации 10 элементов: %v", err)
					return
				}

				style, StyleString, resstyle, err := database.ClothingAdvice(db, StyleSwitchInt) // В данном значении case здесь происходит вызов output ещё два раза, прежде чем получить окончательные данные
				if err != nil {
					log.Printf("Ошибка при выборе города и стиля (case2): %v", err)
					return
				}

				for {
					if resstyle == nil {
						break
					}
					signal := output.PrintClothingAdviceResult(style, StyleString, resstyle) // Печатаем итоговый результат
					if signal == "break" {
						break
					}
				}
			}
		case "q":
			output.Bye() // печатаем финальные слова перед завершением программы
			return
		}
	}
}
