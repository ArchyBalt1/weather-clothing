package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func WeatherFunc() (string, int, error) {
	err := godotenv.Load("F:\\got\\weather-clothing\\.env") // Достали ключ
	if err != nil {
		return "", .0, errors.New("ошибка загрузки из .env файла")
	}

	key := os.Getenv("OPENWEATHER_KEY") // Получили
	var city string
	fmt.Println("Введите название города")
	fmt.Scan(&city)

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, key)

	req, err := http.Get(url)
	if err != nil {
		return "", .0, errors.New("ошибка при получении Get запроса")
	}

	var weather Weather
	err = json.NewDecoder(req.Body).Decode(&weather)
	if err != nil {
		return "", .0, err
	}

	return city, int(math.Round(weather.Main.Temp)), nil
}
