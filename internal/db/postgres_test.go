package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectTestdb() *sql.DB {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		return nil
	}

	constSQL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", constSQL)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		return nil
	}
	return db
}

func TestNotificationConditionsPressureWind_speed(t *testing.T) {
	type args struct {
		conditions string
		pressure   int
		wind_speed float32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Norm", args: args{
			conditions: "Rain",
			pressure:   1010,
			wind_speed: 5.49,
		}, want: fmt.Sprint("На улице дождь — без зонта лучше не выходить 🌦️\nДавление в пределах нормы — хорошее самочувствие ожидается\nЛёгкий ветерок — можно идти в лёгкой куртке")},
	}
	db := ConnectTestdb()
	if db == nil {
		fmt.Println("Подключение к бд с ошибкой")
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotificationConditionsPressureWind_speed(db, tt.args.conditions, tt.args.pressure, tt.args.wind_speed); got != tt.want {
				t.Errorf("NotificationConditionsPressureWind_speed() = %v, want %v", got, tt.want)
			}
		})
	}
	db.Close()
}

func TestWriteWeatherHistory(t *testing.T) {
	type args struct {
		city       string
		temp       int
		conditions string
		pressure   int
		wind_speed float32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "strings",
			args:    args{city: "нОВоСибирСк", temp: 25, conditions: "Clear", pressure: 1000, wind_speed: 5},
			wantErr: false},
	}
	db := ConnectTestdb()
	if db == nil {
		fmt.Println("Подключение к бд с ошибкой")
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteWeatherHistory(db, tt.args.city, tt.args.temp, tt.args.conditions, tt.args.pressure, tt.args.wind_speed); (err != nil) != tt.wantErr {
				t.Errorf("WriteWeatherHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.Close()
}
