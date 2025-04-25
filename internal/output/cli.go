package output

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"weather-clothing/internal/models"
)

func Hello() {
	fmt.Println("📋 Что сегодня хочешь узнать?")
}

func HelloMenu(i *int) {
	if *i > 0 {
		fmt.Println("\nДавай выясним что-то ещё")
	}
	fmt.Print("1: Погодные условия\n2: История запросов погодных условий\n3: Лучшие стили ➕ погодные условия\n> ")
	*i++
}

func WeatherPrint(signal int) {
	if signal == 0 {
		fmt.Print("🌆 Введите название города (или 'q' для выхода в меню):\n> ")
	}
	if signal == 1 {
		fmt.Print("Введён неккоректный город\n> ")
	}
	if signal == 2 {
		fmt.Print("🌆 Давай ещё один (не забывай про 'q'):\n> ")
	}
}

func PrintWeatherResult(city string, temp int, conditions, notification string, wind_speed float32, pressure int) string {
	var y string
	fmt.Printf("📍 %s %d°C, %s\n%s\n\n", city, temp, conditions, notification)
	for {
		fmt.Print("🔍 Хочешь увидеть подробности? (y/n/q)\n> ")
		fmt.Scan(&y)
		if y == "y" || y == "н" {
			fmt.Printf("📊 Подробности:\n• Скорость ветра %.2f м/с\n• Давление %d гПа\nЛюбая клавиша для продолжения...", wind_speed, pressure)
			fmt.Scan(&y)
			fmt.Println()
			break
		} else if y == "n" || y == "т" {
			break
		} else if y == "q" || y == "й" {
			return "break"
		}
	}

	return ""
}

func PrintHistoryRecent_requests(FilterSlice []string) {
	fmt.Println("\n📜 Последние запросы:")
	index := 1
	for _, i := range FilterSlice {
		fmt.Printf("%d: %s\n", index, i)
		index++
	}
}

func PrintHistoryResult(wHistory []models.WeatherHistory_10) string {
	var cityes string
	fmt.Print("Введите название города (или 'q' для выхода в меню):\n> ")
	fmt.Scan(&cityes)

	if cityes == "q" || cityes == "й" {
		return "break"
	}

	j := 1
	for i := 9; i >= 0; i-- {
		if strings.EqualFold(cityes, wHistory[i].City) {
			if j == 1 {
				fmt.Println("\n📋 Недавно запрошенные позиции:")
			}
			fmt.Printf("• %d: %v\n%s %d°C, %s\nWind: %.2f м/c; Pressure: %d гПа\n\n", j, wHistory[i].Date.Format("15:04, 02-01-2006"), wHistory[i].City, wHistory[i].Temp, wHistory[i].Conditions, wHistory[i].Wind_speed, wHistory[i].Pressure)
			j++
		}
	}
	if j == 1 {
		fmt.Println("\n🫨  Введён неккоректный город, давай поднимательнее")
	} else {
		fmt.Println("Может ещё один? 👀")
	}

	return ""
}

func PrintClothingAdviceResult_Hello() {
	fmt.Print("\n🧥 Под какую погоду подоберём стиль?\n1. Последняя запрошенная\n2. Выбрать из 10 последних записей:\n> ")
}

func PrintClothingAdviceResult(style models.Style, StyleString []string, resstyle []models.ResStyle) string {
	fmt.Printf("%s %d°C, %s, %.2fм/с\n", style.City, style.Temp, style.Conditions, style.Wind_speed)
	if StyleString == nil {
		fmt.Println(resstyle[0].Comments)
		return "break"
	}

	NewMap := make(map[int]string)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Выберите стиль (или 'q' для выхода в меню):")
	for index, key := range StyleString {
		fmt.Printf("• %d: %s\n", index+1, key)
		NewMap[index+1] = key
	}
	fmt.Print("> ")
	for {
		var value string
		var b int
		if scanner.Scan() {
			if scanner.Text() == "" {
				continue
			} else if scanner.Text() == "q" || scanner.Text() == "й" {
				return "break"
			} else if scanner.Text() >= "a" && scanner.Text() <= "z" {
				continue
			}

			b, _ = strconv.Atoi(scanner.Text())
			value, _ = NewMap[b]
		}

		IsViewed := false
		for _, i := range resstyle {
			if strings.EqualFold(i.Style, value) {
				fmt.Printf("\n%s:\n%s\n", i.Style, i.Comments)
				fmt.Println()
				IsViewed = true
				delete(NewMap, b)
			}
		}

		if !IsViewed {
			fmt.Println("Такого стиля нет в списке, давай повнимательнее")
		} else if len(NewMap) == 0 {
			fmt.Print("😦 Стили закончились\n")
			return "break"
		} else if IsViewed {
			fmt.Println("Давай посмотрим ещё один стиль 🥷  (или 'q' для выхода в меню)")
		}
		for index, key := range NewMap {
			fmt.Printf("• %d: %s\n", index, key)
		}
		fmt.Print("> ")
	}
}

func PrintClothingAdviceResultHistory(FilterSlice []string, wHistory []models.WeatherHistory_10, style *models.Style) string {
	var a string
	fmt.Println("🏙️  Города и данные:")
	j := 1
	for _, i := range wHistory {
		fmt.Printf("• %d: %v\n%s %d°C, %s\n\n", j, i.Date.Format("15:04, 02-01-2006"), i.City, i.Temp, i.Conditions)
		j++
	}
	fmt.Print("Посмотрим подробности? y/n/q\n> ")
	fmt.Scan(&a)

	if a == "q" || a == "й" {
		return "breakQ"
	}

	if a == "y" || a == "н" {
		for {
			j = 1
			fmt.Print("Введите номер города для просмотра подробностей ('q' для выхода в меню и 's' для продолжения):\n> ")
			fmt.Scan(&a)
			if a == "q" || a == "й" {
				return "breakQ"
			} else if a == "s" || a == "ы" {
				break
			}

			aInt, err := strconv.Atoi(a)
			if err != nil {
				fmt.Println("Введено не число")
				continue
			}
			for _, i := range wHistory {
				if aInt == j {
					fmt.Printf("• %d: %v\n%s %d°C, %s\nWind: %.2f м/c; Pressure: %d гПа\n\n", j, i.Date.Format("15:04, 02-01-2006"), i.City, i.Temp, i.Conditions, i.Wind_speed, i.Pressure)
				}
				j++
			}
		}
	}

	fmt.Print("Выберите желаемый номер города для подборки подходящего стиля (или 'q' для выхода в меню):\n> ")
	j = 1
	for {
		fmt.Scan(&a)
		if a == "q" {
			return "breakQ"
		}

		aInt, err := strconv.Atoi(a)
		if err != nil {
			fmt.Println("Введено не число")
			continue
		}
		if aInt >= 1 && aInt <= 10 {
			for _, i := range wHistory {
				if aInt == j {
					style.City = i.City
					style.Temp = i.Temp
					style.Conditions = i.Conditions
					style.Wind_speed = i.Wind_speed
				}
				j++
			}
			return "break"

		} else {
			fmt.Println("Неккоректный номер")
		}
	}
}
