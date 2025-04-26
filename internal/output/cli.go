package output

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"weather-clothing/internal/models"
)

func Hello() {
	fmt.Println("📋 Выбери пункт меню:")
	fmt.Print("1: Погодные условия\n2: История запросов погодных условий\n3: Лучшие стили ➕ погодные условия\n'q': Завершить программу\n> ")
}

func WeatherPrint(signal int) {
	if signal == 0 {
		fmt.Print("🌆 Введи название города (или 'q' для выхода в меню):\n> ")
	}
	if signal == 1 {
		fmt.Print("Введён неккоректный город\n> ")
	}
	if signal == 2 {
		fmt.Print("🌆 Давай ещё один (не забывай про 'q'):\n> ")
	}
}

func PrintWeatherResult(city string, temp int, conditions, notification string, wind_speed float32, pressure int) string {
	var details string
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("📍 %s %d°C, %s\n%s\n\n", city, temp, conditions, notification)
	for {
		fmt.Print("🔍 Хочешь увидеть подробности? (y/n/q)\n> ")
		details, _ = reader.ReadString('\n')
		details = strings.TrimSpace(details)

		if details == "y" || details == "н" {
			fmt.Printf("📊 Подробности:\n• Скорость ветра %.2f м/с\n• Давление %d гПа\nЛюбая клавиша для продолжения...", wind_speed, pressure)
			reader.ReadString('\n')
			break
		} else if details == "n" || details == "т" {
			break
		} else if details == "q" || details == "й" {
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введи название города (или 'q' для выхода в меню):\n> ")
	city, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Ошибка reader", err)
		return "break"
	}
	city = strings.TrimSpace(city)

	if city == "q" || city == "й" {
		return "break"
	}

	j := 1
	for i := 9; i >= 0; i-- {
		if strings.EqualFold(city, wHistory[i].City) {
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

func PrintClothingAdviceResult(style models.CityStyle, StyleString []string, resstyle []models.ResStyle) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s %d°C, %s, %.2fм/с\n", style.City, style.Temp, style.Conditions, style.Wind_speed)
	if StyleString == nil {
		fmt.Println(resstyle[0].Comments)
		return "break"
	}

	styleMap := make(map[int]string)
	fmt.Println("Выбери стиль(ли) (или 'q' для выхода в меню):")
	for index, key := range StyleString {
		fmt.Printf("• %d: %s\n", index+1, key)
		styleMap[index+1] = key
	}
	fmt.Print("> ")

	var StyleСhoice string
Continue:
	for {
		StyleСhoice, _ = reader.ReadString('\n')
		StyleСhoice = strings.TrimSpace(StyleСhoice)
		if StyleСhoice == "q" || StyleСhoice == "й" {
			return "break"
		}

		runes := []rune(StyleСhoice)
		for i := len(runes) - 1; i >= 0; i-- {
			if !unicode.IsDigit(runes[i]) {
				fmt.Println("🚧 Неверный формат ввода")
				fmt.Print("> ")
				continue Continue
			}
		}

		index, _ := strconv.Atoi(StyleСhoice)
		IsViewed := false
		for _, i := range resstyle {
			if strings.EqualFold(i.Style, styleMap[index]) {
				fmt.Printf("\n%s:\n%s\n🎯 Не забудь взять %s\n", i.Style, i.Comments, strings.ToLower(i.Accessories))
				fmt.Println()
				IsViewed = true
				delete(styleMap, index)
			}
		}

		if !IsViewed {
			fmt.Println("Такого стиля нет в списке, давай повнимательнее")
		} else if len(styleMap) == 0 {
			fmt.Print("😦 Стили закончились\n")
			return "break"
		} else if IsViewed {
			fmt.Println("Давай посмотрим ещё один стиль 🥷  (или 'q' для выхода в меню)")
		}
		for index, key := range styleMap {
			fmt.Printf("• %d: %s\n", index, key)
		}
		fmt.Print("> ")
	}
}

func PrintClothingAdviceResultHistory(wHistory []models.WeatherHistory_10) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("🏙️  Города и данные:")
	j := 1
	for _, i := range wHistory {
		fmt.Printf("• %d: %v\n%s %d°C, %s\n\n", j, i.Date.Format("15:04, 02-01-2006"), i.City, i.Temp, i.Conditions)
		j++
	}

	var StyleDetail string
	for {
		fmt.Print("Посмотрим подробности? y/n/q\n> ")
		StyleDetail, _ = reader.ReadString('\n')
		StyleDetail = strings.TrimSpace(StyleDetail)
		if !(StyleDetail == "y" || StyleDetail == "n" || StyleDetail == "q") {
			fmt.Println("🚧 Неверный формат ввода")
			continue
		}
		break
	}

	if StyleDetail == "q" || StyleDetail == "й" {
		return "breakQ"
	}

	if StyleDetail == "y" || StyleDetail == "н" {
		for {
			j = 1
			fmt.Print("Введи номер города для просмотра подробностей ('q' для выхода в меню и 's' для продолжения):\n> ")
			StyleDetail, _ := reader.ReadString('\n')
			StyleDetail = strings.TrimSpace(StyleDetail)
			if StyleDetail == "q" || StyleDetail == "й" {
				return "breakQ"
			} else if StyleDetail == "s" || StyleDetail == "ы" {
				return ""
			}

			StyleDetailInt, err := strconv.Atoi(StyleDetail)
			if err != nil {
				fmt.Println("🚧 Неверный формат ввода")
				continue
			}
			for _, i := range wHistory {
				if StyleDetailInt == j {
					fmt.Printf("• %d: %v\n%s %d°C, %s\nWind: %.2f м/c; Pressure: %d гПа\n\n", j, i.Date.Format("15:04, 02-01-2006"), i.City, i.Temp, i.Conditions, i.Wind_speed, i.Pressure)
				}
				j++
			}
		}
	}
	return ""
}

func PrintClothingAdviceResultHistoryCity(wHistory []models.WeatherHistory_10, style *models.CityStyle) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Выбери желаемый номер города для подборки подходящего стиля (или 'q' для выхода в меню):\n> ")
	j := 1
	for {
		StyleDetail, _ := reader.ReadString('\n')
		StyleDetail = strings.TrimSpace(StyleDetail)
		if StyleDetail == "q" || StyleDetail == "й" {
			return "breakQ"
		}

		StyleDetailInt, err := strconv.Atoi(StyleDetail)
		if err != nil {
			fmt.Println("🚧 Неверный формат ввода")
			continue
		}
		if StyleDetailInt >= 1 && StyleDetailInt <= 10 {
			for _, i := range wHistory {
				if StyleDetailInt == j {
					style.City = i.City
					style.Temp = i.Temp
					style.Conditions = i.Conditions
					style.Wind_speed = i.Wind_speed
				}
				j++
			}
			return "break"

		} else {
			fmt.Println("🚧 Неверный формат ввода")
			fmt.Print("> ")
		}
	}
}

func Bye() {
	fmt.Println("👋 Возвращайся скорее")
}
