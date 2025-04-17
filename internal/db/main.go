package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	//"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Init() (*sql.DB, error) {
	constSQL := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
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
}

func Opendb(db *sql.DB, city string, temp int) error {
	rows, err := sq.Select("*").From("weather_history").RunWith(db).Query()
	if err != nil {
		log.Println("Ошибка 2")
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var city string
		var temp int
		var date string
		rows.Scan(&id, &city, &temp, &date)
		fmt.Println(id, city, temp, date)
	}

	/*rows, err := sq.Select("*").From("clothing_advice").RunWith(db).Query()
	if err != nil {
		log.Println("Ошибка 2")
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var temp_max int
		var temp_min int
		var style string
		var comments string
		rows.Scan(&id, &comments, &style, &temp_max, &temp_min)
		fmt.Println(id, style, comments, temp_max, temp_min)
	}*/
	return nil
}

func WeatherHistory(db *sql.DB, city string, temp int) error {
	_, err := sq.Insert("weather_history").Columns("city", "temp").Values(city, temp).PlaceholderFormat(sq.Dollar).RunWith(db).Exec() // Добавка в бд
	if err != nil {
		return err
	}

	return nil
}

func ReadHistory(db *sql.DB, city string, temp int) error {
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
	} // Храним только 10 последних запросов */

	rows, err := sq.Select("*").From("weather_history").RunWith(db).Query()
	if err != nil {
		log.Println("Ошибка 2")
		return err
	}
	defer rows.Close()

	i := 1
	for rows.Next() {
		var id int
		var city string
		var temp int
		var date string
		rows.Scan(&id, &city, &temp, &date)
		fmt.Println(i, city, temp, date)
		i++
	}

	return nil
}
