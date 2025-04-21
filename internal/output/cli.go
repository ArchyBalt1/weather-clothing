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
	fmt.Println("Выберите пункт\n1: Погодные условия\n2: История запросов (последние 10 записей)\n3: Совет по luck(ограниченный доступ)")
}

func PrintWeatherResult(city string, temp int, conditions, notification string, wind_speed float32, pressure int) {
	var y string
	fmt.Printf("%s %d°C, %s\n%s\n\n", city, temp, conditions, notification)
	for {
		fmt.Print("Вывести подробности? (y/n)\n> ")
		fmt.Scan(&y)
		if y == "y" || y == "н" {
			fmt.Printf("📊 Подробности:\n• Скорость ветра %.2f м/с\n• Давление %d гПа\nЛюбая клавиша для продолжения...", wind_speed, pressure)
			fmt.Scan(&y)
			break
		} else if y == "n" || y == "т" {
			break
		}
	}
}

func PrintHistoryResult(FilterSlice []string, cityes string, wHistory []models.WeatherHistory_10) string {
	fmt.Println("Недавно запрошенные города:")
	for _, i := range FilterSlice {
		fmt.Printf("> %s\n", i)
	}
	fmt.Println("Введите название для просмотра погоды или q для возврата в меню:")
	fmt.Scan(&cityes)

	if cityes == "q" || cityes == "й" {
		return "break"
	}

	j := 1
	for i := 9; i >= 0; i-- {
		if cityes == wHistory[i].City {
			fmt.Printf("Number: %d\n• %s %d°C, %s; Wind: %.2f м/c; Pressure: %d гПа; Time: %v\n", j, wHistory[i].City, wHistory[i].Temp, wHistory[i].Conditions, wHistory[i].Wind_speed, wHistory[i].Pressure, wHistory[i].Date.Format("15:04:05 02-01-2006"))
			j++
		}
	}
	if j == 1 {
		fmt.Println("Введён неккоректный город, давай повнмательнее")
	} else {
		fmt.Println("Может ещё один?")
	}

	return ""
}

func PrintClothingAdviceResult_Hello() {
	fmt.Print("Под какую погоду подоберём стиль?\n1. Последняя запрошенная\n2. Выбрать из 10 последних записей:\n")
}

func PrintClothingAdviceResult_LastEntry(style models.Style, StyleString []string, resstyle []models.ResStyle) string {
	fmt.Printf("%s %d°C, %s, %.2fм/с\n", style.City, style.Temp, style.Conditions, style.Wind_speed)
	//fmt.Println("2", StyleString)
	//fmt.Println("сигнал")
	if StyleString == nil {
		fmt.Println(resstyle[0].Comments)
		return "break"
	}

	NewMap := make(map[int]string)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Выберите стиль или нажмите q для выхода в меню:")
	for index, key := range StyleString {
		fmt.Printf("• %d: %s\n", index+1, key)
		NewMap[index+1] = key
	}
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
				fmt.Printf("%s:\n%s\n", i.Style, i.Comments)
				fmt.Println()
				IsViewed = true
				delete(NewMap, b)
			}
		}

		if !IsViewed {
			fmt.Println("Такого стиля нет в списке, давай повнимательнее")
		} else if len(NewMap) == 0 {
			fmt.Print("Стили закончились\n\n")
			return "break"
		} else if IsViewed {
			fmt.Println("Хочешь посмотреть другой стиль?")
		}
		for index, key := range NewMap {
			fmt.Printf("• %d: %s\n", index, key)
		}
	}
}
