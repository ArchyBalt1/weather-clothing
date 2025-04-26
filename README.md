# weather-clothing

__weather-clothing__ - pet проект, который позволяет вам узнать _актуальную погоду_ по выбранному городу/селу/району, _даёт советы по стилю_ одежды в зависимости от погодных условий, а также _сохраняет историю запросов в базу данных_

В данном проекте представлена как CLI-версия (изначально разрабатывалось под неё), а также TelegramBot с 85% функционалом CLI версии.
# Установка
* Клонируйте репозиторий

```bash
https://github.com/ArchyBalt1/weather-clothing.git
```
* Настройте переменные окружения с помощью папки __.env__
  
  Создайте файл .env и укажите:
```ini
OPENWEATHER_KEY=your_token_here
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
TGBOTTOKEN=your_telegram_token
```
_Файл .env.test использовался для теста бд, где была создана отдельная тестовая база данных. [API](https://openweathermap.org/api) для получения погодных условий. Чтобы получить tgbot_token, воспользуйтесь **@BotFather** в телеграмм_
* Настройте базу данных
  
  В корне проекта есть файл **dump.sql**, в котором полностью готовая к работе база данных, со схемой и данными. Для работы с ней:
    * Создайте пустую базу в PostgresSQL
    * Выполните команду
```bash
psql -U your_user -d your_dbname -f dump.sql
```
  * Установите зависимости
```bash
go mod tidy
```
# Запуск
* Для CLI версии

  Перейдите в папку **server** (weather-clothing/cmd/server) и запустите:
```bash
go run main.go
```
* Для Telegram бота
  
  При запуске CLI версии будет выбор, что запустить
```os.Stdin
Хотите запустить TelegramBot? y/any_key
>
```
