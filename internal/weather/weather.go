package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
)

type WeatherPlanet struct {
	Weather []struct {
		Main string `json:"main"` // состояние
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`     // температура
		Pressure int     `json:"pressure"` // давление
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"` // ветер
	}
}

func WeatherFunc(city string) (string, int, string, int, float32, error) {
	key := os.Getenv("OPENWEATHER_KEY") // Получили

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, key)

	req, err := http.Get(url)
	if err != nil {
		return "", .0, "", .0, .0, err
	}

	var weather WeatherPlanet
	err = json.NewDecoder(req.Body).Decode(&weather)
	if err != nil {
		return "", .0, "", .0, .0, err
	}

	if weather.Main.Pressure == 0 {
		return "Введён неккоректный город", .0, "", .0, .0, nil
	}

	return city, int(math.Round(weather.Main.Temp)), weather.Weather[0].Main, weather.Main.Pressure, weather.Wind.Speed, nil

}
