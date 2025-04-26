package logic

import (
	"log"
	"os"
	"time"
	"weather-clothing/internal/models"
)

func LogFile() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0123)
	if err != nil {
		log.Fatalln("Ошибка при открытии файла логов", err)
	}

	log.SetOutput(file)

	log.SetFlags(log.Ldate | log.Ltime)
} // Файл логов

func FilterMap(Slicecity []string, wHistory []models.WeatherHistory_10) []string {
	FilterMap := make(map[string]struct{})
	for _, i := range Slicecity {
		FilterMap[i] = struct{}{}
	}

	FilterSlice := make([]string, 0, 10)
	for _, i := range Slicecity {
		if _, ok := FilterMap[i]; ok {
			FilterSlice = append(FilterSlice, i)
			delete(FilterMap, i)
		}
	}

	for i, j := 0, len(FilterSlice)-1; i < j; i, j = i+1, j-1 {
		FilterSlice[i], FilterSlice[j] = FilterSlice[j], FilterSlice[i]
	}

	return FilterSlice
} // Мапа для фильтрации городов и их вывод в формате первый пришёл-первый в списке

func TimeMonth() string {
	month := time.Now().Month()
	switch month {
	case 12, 1, 2:
		return "Зима"
	case 3, 4, 5:
		return "Весна"
	case 6, 7, 8:
		return "Лето"
	case 9, 10, 11:
		return "Осень"
	}

	return "Любое"
} // Выборка времени года, для более чёткого подбора стиля
