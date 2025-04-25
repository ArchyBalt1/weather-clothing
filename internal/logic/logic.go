package logic

import (
	"log"
	"os"
	"weather-clothing/internal/models"
)

func LogFile() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0123)
	if err != nil {
		log.Fatalln("Ошибка при открытии файла логов", err)
	}

	log.SetOutput(file)

	log.SetFlags(log.Ldate | log.Ltime)
}

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
}

/*func Filter(Slicecity []string, wHistory []models.WeatherHistory_10, signals int) {
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

	for {
		signal := output.PrintHistoryResult(FilterSlice, wHistory)
		if signal == "break" {
			break
		}
	}
	wHistory = nil
	Slicecity = nil
	FilterMap = nil
} // мапа как фильтр */
