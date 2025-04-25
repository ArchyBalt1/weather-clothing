package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"weather-clothing/internal/logic"
	"weather-clothing/internal/models"
	"weather-clothing/internal/output"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

func Init() (*sql.DB, error) {
	constSQL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", constSQL)
	if err != nil {
		log.Println("Ошибка при запуске бд")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Ошибка при запуске бд")
		return nil, err
	}
	return db, nil
} // запуск бд

func WriteWeatherHistory(db *sql.DB, city string, temp int, conditions string, pressure int, wind_speed float32) error {
	_, err := sq.Insert("weather_history").Columns("city", "temp", "conditions", "pressure", "wind_speed").Values(city, temp, conditions, pressure, wind_speed).PlaceholderFormat(sq.Dollar).RunWith(db).Exec() // Добавка в бд (сохранение в формате Qwer)
	if err != nil {
		return err
	}

	return nil
} // insert запрос

func HistoryLimit10(db *sql.DB) error {
	sql1, args, err := sq.
		Select("id").
		From("weather_history").
		OrderBy("id DESC").
		Limit(10).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = sq.
		Delete("weather_history").
		Where(fmt.Sprintf("id NOT IN (%s)", sql1), args...).
		RunWith(db).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		log.Println("Ошибка при delete запросе")
		return err
	} // Храним только 10 последних запросов

	return nil
}

func ReadHistory(db *sql.DB) ([]string, []models.WeatherHistory_10, error) {
	rows, err := sq.Select("city", "temp", "conditions", "pressure", "wind_speed", "date").From("weather_history").RunWith(db).Query()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var wHistory []models.WeatherHistory_10

	Slicecity := make([]string, 0, 10)
	for rows.Next() {
		var w models.WeatherHistory_10
		if err := rows.Scan(&w.City, &w.Temp, &w.Conditions, &w.Pressure, &w.Wind_speed, &w.Date); err != nil {
			return nil, nil, err
		}

		wHistory = append(wHistory, w)
		Slicecity = append(Slicecity, w.City)
	}

	return Slicecity, wHistory, nil
} // работа с историей загрузкит в бд (10 записей. Можем выбрать, какой город посмотреть)

func ClothingAdvice(db *sql.DB, b int) (models.Style, []string, []models.ResStyle, error) {
	switch b {
	case 1:
		row := sq.Select("city", "temp", "conditions", "wind_speed").From("weather_history").OrderBy("id DESC").Limit(1).PlaceholderFormat(sq.Dollar).RunWith(db).QueryRow()

		var style models.Style
		row.Scan(&style.City, &style.Temp, &style.Conditions, &style.Wind_speed)

		var resstyle []models.ResStyle
		StyleString, err := Advice(db, style, &resstyle)
		if err != nil {
			return style, nil, nil, err
		}

		return style, StyleString, resstyle, nil
	case 2:
		Slicecity, wHistory, err := ReadHistory(db)
		if err != nil {
			return models.Style{}, nil, nil, err
		}
		FilterSlice := logic.FilterMap(Slicecity, wHistory)

		var style models.Style
		var resstyle []models.ResStyle
		StyleString := []string{}
		for {
			signal := output.PrintClothingAdviceResultHistory(FilterSlice, wHistory, &style)
			if signal == "breakQ" {
				return models.Style{}, nil, nil, nil
			}

			StyleString, err = Advice(db, style, &resstyle)
			if err != nil {
				return models.Style{}, nil, nil, err
			}
			if signal == "break" {
				break
			}
		}

		return style, StyleString, resstyle, nil
	}
	return models.Style{}, nil, nil, nil
} // стиль

func ClothingAdviceHistory(db *sql.DB, style models.Style) ([]string, []models.ResStyle, error) {
	var resstyle []models.ResStyle
	StyleString, err := Advice(db, style, &resstyle)
	if err != nil {
		return nil, nil, err
	}

	return StyleString, resstyle, nil
}

func Advice(db *sql.DB, style models.Style, resstyle *[]models.ResStyle) ([]string, error) {
	rows, err := sq.Select("style", "comments").From("clothing_advice").Where(sq.And{
		sq.LtOrEq{"temp_min": style.Temp},
		sq.GtOrEq{"temp_max": style.Temp},
	}).Where(" ? = ANY(conditions)", style.Conditions).Where(sq.GtOrEq{"max_speed": style.Wind_speed}).PlaceholderFormat(sq.Dollar).RunWith(db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var StyleString []string
	for rows.Next() {
		var rs models.ResStyle
		err := rows.Scan(&rs.Style, &rs.Comments)
		if err != nil {
			return nil, err
		}
		*resstyle = append(*resstyle, rs)
	}
	//fmt.Println(*resstyle)

	if len(*resstyle) == 1 && (*resstyle)[0].Style == "Pop" {
		return nil, nil
	} else {
		for _, i := range *resstyle {
			if i.Style == "Pop" {
				continue
			}
			StyleString = append(StyleString, i.Style)
		}
	}
	//fmt.Println("1", StyleString)
	return StyleString, nil
} // Поиск и вывод стилей

func NotificationConditionsPressureWind_speed(db *sql.DB, conditions string, pressure int, wind_speed float32) string {
	var wg sync.WaitGroup
	var conditionsComments string
	var pressureComments string
	var wind_speedComments string
	wg.Add(3)
	go func() {
		row := db.QueryRow(`SELECT conditions_comments FROM conditions_advice WHERE conditions = $1 ORDER BY Random() LIMIT(1)`, conditions)
		row.Scan(&conditionsComments)
		wg.Done()
	}()
	go func() {
		row := db.QueryRow(`SELECT pressure_comments FROM pressure_advice WHERE $1 BETWEEN min_pressure AND max_pressure ORDER BY Random() LIMIT(1)`, pressure)
		row.Scan(&pressureComments)
		wg.Done()
	}()
	go func() {
		row := db.QueryRow(`SELECT wind_comments FROM wind_advice WHERE $1 BETWEEN min_speed AND max_speed ORDER BY Random() LIMIT(1)`, wind_speed)
		row.Scan(&wind_speedComments)
		wg.Done()
	}()
	wg.Wait()

	return fmt.Sprintf("%s\n%s\n%s", conditionsComments, pressureComments, wind_speedComments)
} // Советы
