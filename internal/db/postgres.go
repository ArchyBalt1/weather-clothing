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
	_, err := sq.Insert("weather_history").Columns("city", "temp", "conditions", "pressure", "wind_speed").Values(city, temp, conditions, pressure, wind_speed).PlaceholderFormat(sq.Dollar).RunWith(db).Exec() // Добавка в бд
	if err != nil {
		return err
	}

	return nil
} // insert запрос

func ReadHistory(db *sql.DB) error {
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

	rows, err := sq.Select("city", "temp", "conditions", "pressure", "wind_speed", "date").From("weather_history").RunWith(db).Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	var wHistory []models.WeatherHistory_10

	var cityes string
	Slicecity := make([]string, 0, 10)
	for rows.Next() {
		var w models.WeatherHistory_10
		if err := rows.Scan(&w.City, &w.Temp, &w.Conditions, &w.Pressure, &w.Wind_speed, &w.Date); err != nil {
			return err
		}

		wHistory = append(wHistory, w)
		Slicecity = append(Slicecity, w.City)
	}

	for {
		FilterSlice := logic.Filter(Slicecity)
		signal := output.PrintHistoryResult(FilterSlice, cityes, wHistory)
		if signal == "break" {
			break
		}
	} // работа с историей загрузкит в бд (10 записей. Можем выбрать, какой город посмотреть)
	return nil
}

func ClothingAdvice(db *sql.DB) error {
	var a int
	output.PrintClothingAdviceResult_Hello()
	fmt.Scan(&a)
	switch a {
	case 1:
		row := sq.Select("city", "temp", "conditions", "wind_speed").From("weather_history").OrderBy("id DESC").Limit(1).PlaceholderFormat(sq.Dollar).RunWith(db).QueryRow()

		var style models.Style
		row.Scan(&style.City, &style.Temp, &style.Conditions, &style.Wind_speed)

		var resstyle []models.ResStyle
		StyleString, err := Advice(db, style, &resstyle)
		if err != nil {
			return err
		}

		for {
			signal := output.PrintClothingAdviceResult_LastEntry(style, StyleString, resstyle)
			if signal == "break" {
				break
			}
		}

	}
	return nil
} // стиль

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
