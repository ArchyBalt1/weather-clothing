package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	database "weather-clothing/internal/db"
	"weather-clothing/internal/weather"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Привет, я подскажу тебе погоду☀️🌧, а также предложу стили👗🧥")
			bot.Send(msg)
		case "/weather":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Посмотри погоду, введи название города")
			bot.Send(msg)

			cityesUpdate := waitForUserResponse(bot, update.Message.Chat.ID)
			if cityesUpdate == nil {
				msg = tgbot.NewMessage(update.Message.Chat.ID, "Не получили название города, попробуйте снова.")
				bot.Send(msg)
				return
			}

			cityes := cityesUpdate.Message.Text

			city, temp, conditions, pressure, wind_speed, err := weather.WeatherFunc(cityes)
			if err != nil {
				log.Println("Ошибка при получении погодных условий:", err)
				msg = tgbot.NewMessage(update.Message.Chat.ID, "Не удалось получить данные о погоде. Попробуйте другой город.")
				bot.Send(msg)
				return
			}

			// Сохраняем в базу данных
			if err := database.WriteWeatherHistory(db, city, temp, conditions, pressure, wind_speed); err != nil {
				log.Println("Ошибка при сохранении в базу данных:", err)
			}

			// Проверяем условия для уведомлений
			notification := database.NotificationConditionsPressureWind_speed(db, conditions, pressure, wind_speed)

			// Формируем основное сообщение с погодой
			weatherMsg := fmt.Sprintf("🌤 Погода в %s:\nТемпература: %d°C\nСостояние: %s\n%s",
				city, temp, conditions, notification)

			// Создаем inline-клавиатуру для подробностей
			detailedInfo := fmt.Sprintf("📊 Подробности:\n• Скорость ветра: %.2f м/с\n• Давление: %d гПа",
				wind_speed, pressure)

			msg = tgbot.NewMessage(update.Message.Chat.ID, weatherMsg)
			msg.ReplyMarkup = tgbot.NewInlineKeyboardMarkup(
				tgbot.NewInlineKeyboardRow(
					tgbot.NewInlineKeyboardButtonData("Подробнее", "show_details"),
				),
			)

			if _, err := bot.Send(msg); err != nil {
				log.Println("Ошибка отправки сообщения:", err)
			}

			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "show_details" {
					// Отправляем подробную информацию
					detailedMsg := tgbot.NewMessage(
						update.CallbackQuery.Message.Chat.ID,
						detailedInfo, // Переменная из предыдущего кода
					)

					// Удаляем inline-клавиатуру
					edit := tgbot.NewEditMessageReplyMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						tgbot.InlineKeyboardMarkup{InlineKeyboard: [][]tgbot.InlineKeyboardButton{}},
					)

					bot.Send(edit)
					bot.Send(detailedMsg)

					// Подтверждаем обработку callback
					//bot.AnswerCallbackQuery(tgbot.NewCallback(update.CallbackQuery.ID, ""))
				}
			}
		case "/help":
			msg := tgbot.NewMessage(update.Message.Chat.ID, "/start, /weather")
			bot.Send(msg)
		}
	}
}

func waitForUserResponse(bot *tgbot.BotAPI, chatID int64) *tgbot.Update {
	// Таймаут ожидания - 60 секунд
	timeout := time.Now().Add(60 * time.Second)

	for time.Now().Before(timeout) {
		updates, err := bot.GetUpdates(tgbot.NewUpdate(0))
		if err != nil {
			log.Println("Ошибка получения обновлений:", err)
			continue
		}

		for _, update := range updates {
			if update.Message != nil && update.Message.Chat.ID == chatID {
				return &update
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
