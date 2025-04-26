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

var weatherstate = make(map[int64]string)            // отслеживание состояния для переключения режимов
var weathercallback = make(map[int64]string)         // отслеживание кнопок
var styleMap = make(map[int]string)                  // кеш стилей
var style_choice_city = make(map[int64]models.Style) // кеш города из выборки в /style

const menu = `🏠 Главное меню:
🔹 /weather — Узнать текущую погоду
🔹 /history — История последних запросов
🔹 /style — Подобрать стиль по погоде
🔹 /h — Справка по командам`

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
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			switch {
			case data == "show_details":
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

			case strings.HasPrefix(data, "history_"):
				_, wHistory, err := database.ReadHistory(db)
				if err != nil {
					log.Println(err)
					return
				}

				city := strings.TrimPrefix(data, "history_")
				var msgHistory string
				msgHistory = fmt.Sprintln("📜 Недавно запрошенная(ые) позиция(ии)")

				j := 1
				for i := len(wHistory) - 1; i >= 0; i-- {
					if strings.EqualFold(wHistory[i].City, city) {
						msgHistory += fmt.Sprintf("• %d: %v\n%s %d°C, %s\nWind: %.2f м/c; Pressure: %d гПа\n\n", j, wHistory[i].Date.Format("15:04, 02-01-2006"), wHistory[i].City, wHistory[i].Temp, wHistory[i].Conditions, wHistory[i].Wind_speed, wHistory[i].Pressure)
						j++
					}
				}

				if j == 1 {
					msgHistory = "История не найдена. Попробуй другой город."
				}

				bot.Send(tgbot.NewMessage(chatID, msgHistory))
				bot.Request(tgbot.NewCallback(update.CallbackQuery.ID, ""))
			}

			continue
		}

		text := strings.TrimSpace(strings.ToLower(update.Message.Text))
		chatID := update.Message.Chat.ID

		if strings.EqualFold(text, "назад") {
			weatherstate[chatID] = ""
			bot.Send(tgbot.NewMessage(chatID, menu))
			continue
		}

		switch text {
		case "/start":
			weatherstate[chatID] = ""
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Привет, я подскажу тебе погоду☀️🌧, а также предложу стили👗🧥\n/h: Справка по командам\n/m: Меню")
			bot.Send(msg)
		case "/m":
			weatherstate[chatID] = ""
			bot.Send(tgbot.NewMessage(chatID, menu))
		case "/weather":
			weatherstate[chatID] = "city_request"
			msg := tgbot.NewMessage(chatID, "Посмотри погоду, введи название города (или 'назад' для возврата в меню)")
			bot.Send(msg)
			continue
		case "/history":
			weatherstate[chatID] = "history_request"
		case "/style":
			weatherstate[chatID] = "style_request_city"
			msg := tgbot.NewMessage(chatID, "Под какую погоду подоберём стиль?\n1. Последняя запрошенная\n2. Выбрать из 10 последних записей:")
			bot.Send(msg)
			continue
		case "/h":
			weatherstate[chatID] = ""
			TextHelp := `/start: С чего всё начинать, приветствие
/weather: Текущая погода в любой точке мира + краткие реплики по данным погоды
/history: Недавно запрошенная погода, последние 10 записей
/style: Подберём стиль либо под последнюю запрошенную погоду, либо под недавно запрошенную
"назад": Возврат в меню`
			msg := tgbot.NewMessage(update.Message.Chat.ID, TextHelp)
			bot.Send(msg)
		}

		switch weatherstate[chatID] {
		case "city_request":
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

			notification := database.NotificationConditionsPressureWind_speed(db, temp, conditions, pressure, wind_speed)

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

			buttons := make([][]tgbot.InlineKeyboardButton, 0)
			for _, city := range FilterSlice {
				btn := tgbot.NewInlineKeyboardButtonData(city, "history_"+strings.ToLower(city))
				buttons = append(buttons, tgbot.NewInlineKeyboardRow(btn))
			}
			keyboard := tgbot.NewInlineKeyboardMarkup(buttons...)

			msg := tgbot.NewMessage(chatID, "Выбери город, чтобы посмотреть историю:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "style_choice_city":
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
			} else {
				log.Println("Неккоректный номер")
			}
			style_choice_city[chatID] = style
			StyleString, resstyle, err := database.ClothingAdviceHistory(db, style)
			if err != nil {
				log.Println("Ошибка при запросе бд из тг", err)
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
				styleMap[index+1] = key
			}
			bot.Send(tgbot.NewMessage(chatID, msgStyle))

			weatherstate[chatID] = "style_request_style"

		case "style_request_city":
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
					MsgCityNumber += fmt.Sprintf("• %d: 🕒 %v\n%s %d°C, %s\n\n", j, i.Date.Format("15:04, 01-02-2006"), i.City, i.Temp, i.Conditions)
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
				styleMap[index+1] = key
			}
			bot.Send(tgbot.NewMessage(chatID, msgStyle))
			weatherstate[chatID] = "style_request_style"

		case "style_request_style":
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
					log.Println("Ошибка при запросе бд из тг", err)
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
				value, _ := styleMap[TextInt]
				if value == i.Style {
					msgResStyle = fmt.Sprintf("%s:\n%s\n", i.Style, i.Comments)
					IsViewed = true
					delete(styleMap, TextInt)
				}
			}
			bot.Send(tgbot.NewMessage(chatID, msgResStyle))

			if !IsViewed {
				bot.Send(tgbot.NewMessage(chatID, "Такого стиля нет в списке, давай повнимательнее"))
				continue
			} else if len(styleMap) == 0 {
				weatherstate[chatID] = ""
				bot.Send(tgbot.NewMessage(chatID, "Стили закончились\nВозврат в главное меню"))
				continue
			} else if IsViewed {
				msgStyle := fmt.Sprintln("Хочешь посмотреть другой стиль?")
				for key, value := range styleMap {
					msgStyle += fmt.Sprintf("• %d: %s\n", key, value)
				}
				bot.Send(tgbot.NewMessage(chatID, msgStyle))
			}
		}
	}
}
