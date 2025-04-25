package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	database "weather-clothing/internal/db"
	"weather-clothing/internal/logic"
	"weather-clothing/internal/models"
	"weather-clothing/internal/weather"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var weatherstate = make(map[int64]string)
var weathercallback = make(map[int64]string)
var NewMap = make(map[int]string)
var style_choice_city = make(map[int64]models.Style)

func Bot(db *sql.DB) {
	tgtoken := os.Getenv("TGBOTTOKEN")
	if tgtoken == "" {
		log.Println("Ошибка получения токена")
		return
	}

	bot, err := tgbot.NewBotAPI(tgtoken)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Авторизация аккаунта %s", bot.Self.UserName)

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Сначала проверяем callback
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			if update.CallbackQuery.Data == "show_details" {
				detailinfo := weathercallback[chatID]

				detailedMsg := tgbot.NewMessage(chatID, detailinfo)

				edit := tgbot.NewEditMessageReplyMarkup(
					chatID,
					update.CallbackQuery.Message.MessageID,
					tgbot.InlineKeyboardMarkup{InlineKeyboard: [][]tgbot.InlineKeyboardButton{}},
				)

				bot.Send(edit)
				bot.Send(detailedMsg)

				bot.Request(tgbot.NewCallback(update.CallbackQuery.ID, "Данные загружены 👍"))
			}
			continue
		}

		text := strings.TrimSpace(strings.ToLower(update.Message.Text))
		chatID := update.Message.Chat.ID

		switch weatherstate[chatID] {
		case "city_request":
			if text == "назад" {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Возврат в главное меню"))
				continue
			}

			city, temp, conditions, pressure, wind_speed, err := weather.WeatherFunc(text)
			if city == "Введён неккоректный город" {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "Введён неккоректный город, попробуй ещё")
				bot.Send(msg)
				continue
			}
			if err != nil {
				log.Println("Ошибка при получении погодных условий:", err)
				weatherstate[chatID] = ""
				continue
			}

			if err := database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed); err != nil {
				log.Println("Ошибка при сохранении в базу данных:", err)
			}

			notification := database.NotificationConditionsPressureWind_speed(db, conditions, pressure, wind_speed)

			weatherMsg := fmt.Sprintf("🌤 %s %d°C, %s\n%s\n\n", city, temp, conditions, notification)

			detailedInfoMsg := fmt.Sprintf("📊 Подробности:\n• Скорость ветра: %.2f м/с\n• Давление: %d гПа", wind_speed, pressure)

			weathercallback[chatID] = detailedInfoMsg

			msg := tgbot.NewMessage(chatID, weatherMsg)
			msg.ReplyMarkup = tgbot.NewInlineKeyboardMarkup(
				tgbot.NewInlineKeyboardRow(
					tgbot.NewInlineKeyboardButtonData("Подробнее", "show_details"),
				),
			)

			bot.Send(msg)
		case "history_request":
			if text == "назад" || text == "2" {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Возврат в главное меню"))
				continue
			}
			if err := database.HistoryLimit10(db); err != nil {
				log.Println(err)
				return
			} // фильтруем 10 последних записей
			Slicecity, wHistory, err := database.ReadHistory(db)
			if err != nil {
				log.Println(err)
				return
			} // логика выборки
			FilterSlice := logic.FilterMap(Slicecity, wHistory)

			if text == "список" || text == "1" {
				var msgCity string
				for _, i := range FilterSlice {
					msgCity += fmt.Sprintf("> %s\n", i)
				}
				bot.Send(tgbot.NewMessage(chatID, msgCity))
				continue
			}

			signal := 1
			var msgHistory string
			msgHistory = fmt.Sprintln("Недавно запрошенные позиции:")
			for i := 9; i >= 0; i-- {
				if text == strings.ToLower(wHistory[i].City) {
					msgHistory += fmt.Sprintf("• %s %d°C, %s; Wind: %.2f м/c; Pressure: %d гПа; Time: %v\n", wHistory[i].City, wHistory[i].Temp, wHistory[i].Conditions, wHistory[i].Wind_speed, wHistory[i].Pressure, wHistory[i].Date.Format("15:04:05 02-01-2006"))
					signal++
				}
			}
			bot.Send(tgbot.NewMessage(chatID, msgHistory))
			if signal == 1 {
				bot.Send(tgbot.NewMessage(chatID, "Введён неккоректный город, давай повнимательнее"))
				continue
			} else {
				bot.Send(tgbot.NewMessage(chatID, "Может ещё один?"))
			}

		case "style_choice_city":
			if text == "назад" {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Возврат в главное меню"))
				continue
			}

			_, wHistory, err := database.ReadHistory(db)
			if err != nil {
				continue
			}

			var style models.Style
			j := 1
			textInt, _ := strconv.Atoi(text)
			if textInt >= 1 && textInt <= 10 {
				for _, i := range wHistory {
					if textInt == j {
						style.City = i.City
						style.Temp = i.Temp
						style.Conditions = i.Conditions
						style.Wind_speed = i.Wind_speed
					}
					j++
				}
				fmt.Println(style)
			} else {
				fmt.Println("Неккоректный номер")
			}
			style_choice_city[chatID] = style
			StyleString, resstyle, err := database.ClothingAdviceHistory(db, style)
			if err != nil {
				log.Println("Arpol")
			}

			if StyleString == nil {
				msgPop := fmt.Sprint(resstyle[0].Comments)
				bot.Send(tgbot.NewMessage(chatID, msgPop))
				weatherstate[chatID] = ""
				continue
			}

			msgWeatherCity := tgbot.NewMessage(chatID, fmt.Sprintf("%s %d°C, %s, %.2fм/с\n", style.City, style.Temp, style.Conditions, style.Wind_speed))
			bot.Send(msgWeatherCity)

			msgStyle := fmt.Sprintln("Выберите стиль или 'назад' для возврата в меню:")
			for index, key := range StyleString {
				msgStyle += fmt.Sprintf("• %d: %s\n", index+1, key)
				NewMap[index+1] = key
			}
			bot.Send(tgbot.NewMessage(chatID, msgStyle))

			weatherstate[chatID] = "style_request_style"

		case "style_request_city":
			if text == "назад" {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Возврат в главное меню"))
				continue
			}
			textInt, _ := strconv.Atoi(text)
			if textInt == 2 {
				if err := database.HistoryLimit10(db); err != nil {
					log.Println(err)
					return
				} // фильтруем 10 последних записей

				bot.Send(tgbot.NewMessage(chatID, "Города и данные, введи номер города:"))
				_, wHistory, err := database.ReadHistory(db)
				if err != nil {
					continue
				}

				j := 1
				var MsgCityNumber string
				for _, i := range wHistory {
					MsgCityNumber += fmt.Sprintf("Number: %d\n• %s %d°C, %s %v\n", j, i.City, i.Temp, i.Conditions, i.Date.Format("15:04 01-02-2006"))
					j++
				}

				bot.Send(tgbot.NewMessage(chatID, MsgCityNumber))
				weatherstate[chatID] = "style_choice_city"
				continue
			}

			style, StyleString, resstyle, err := database.ClothingAdvice(db, textInt)
			if err != nil {
				log.Println(err)
				return
			}

			msgWeatherCity := tgbot.NewMessage(chatID, fmt.Sprintf("%s %d°C, %s, %.2fм/с\n", style.City, style.Temp, style.Conditions, style.Wind_speed))
			bot.Send(msgWeatherCity)

			if StyleString == nil {
				msgPop := fmt.Sprint(resstyle[0].Comments)
				bot.Send(tgbot.NewMessage(chatID, msgPop))
				weatherstate[chatID] = ""
				continue
			}

			msgStyle := fmt.Sprintln("Выберите стиль или 'назад' для возврата в меню:")
			for index, key := range StyleString {
				msgStyle += fmt.Sprintf("• %d: %s\n", index+1, key)
				NewMap[index+1] = key
			}
			bot.Send(tgbot.NewMessage(chatID, msgStyle))
			weatherstate[chatID] = "style_request_style"
			fmt.Println(NewMap) //

		case "style_request_style":
			if text == "назад" {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Возврат в главное меню"))
				continue
			}

			TextInt, err := strconv.Atoi(text)
			if err != nil {
				log.Println("Неверное преобразоваение")
				weatherstate[chatID] = ""
				continue
			}

			var resstyle []models.ResStyle
			value, ok := style_choice_city[chatID]
			if ok {
				_, resstyle, err = database.ClothingAdviceHistory(db, value)
				if err != nil {
					log.Println("Arpol")
				}
			} else {
				_, _, resstyle, err = database.ClothingAdvice(db, 1)
				if err != nil {
					log.Println(err)
					return
				}
			}

			IsViewed := false
			var msgResStyle string
			for _, i := range resstyle {
				value, _ := NewMap[TextInt]
				if value == i.Style {
					msgResStyle = fmt.Sprintf("%s:\n%s\n", i.Style, i.Comments)
					IsViewed = true
					delete(NewMap, TextInt)
				}
			}
			bot.Send(tgbot.NewMessage(chatID, msgResStyle))

			if !IsViewed {
				bot.Send(tgbot.NewMessage(chatID, "Такого стиля нет в списке, давай повнимательнее"))
				continue
			} else if len(NewMap) == 0 {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Стили закончились\nВозврат в главное меню"))
				continue
			} else if IsViewed {
				msgStyle := fmt.Sprintln("Хочешь посмотреть другой стиль?")
				for key, value := range NewMap {
					msgStyle += fmt.Sprintf("• %d: %s\n", key, value)
				}
				bot.Send(tgbot.NewMessage(chatID, msgStyle))
			}
		}
		switch text {
		case "/start":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Привет, я подскажу тебе погоду☀️🌧, а также предложу стили👗🧥\n/h")
			bot.Send(msg)
		case "/weather":
			weatherstate[chatID] = "city_request"
			msg := tgbot.NewMessage(chatID, "Посмотри погоду, введи название города (или 'назад' для возврата в меню)")
			bot.Send(msg)
		case "/history":
			weatherstate[chatID] = "history_request"
			msg := tgbot.NewMessage(chatID, "1: 'Список'\n2: 'Назад'\nНазвание города")
			bot.Send(msg)
		case "/style":
			weatherstate[chatID] = "style_request_city"
			msg := tgbot.NewMessage(chatID, "Под какую погоду подоберём стиль?\n1. Последняя запрошенная\n2. Выбрать из 10 последних записей:")
			bot.Send(msg)
		case "/h":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "/start, /weather, /history, /style")
			bot.Send(msg)
		}
	}
}
