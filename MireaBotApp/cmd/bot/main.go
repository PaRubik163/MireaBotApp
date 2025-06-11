package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	handler "mireabot/internal/parser/bot"
	database "mireabot/internal/parser/bot/storage"
)

type UserState struct {
	login            string
	password         string
	awaitingLogin    bool
	awaitingPassword bool
	updateLogin      bool
	uppdatePassword  bool
}

var userStates = make(map[int64]*UserState)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("Short BotToken")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	//bot.Debug = true
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		// Если пользователя нет в мапе — создаём запись
		if _, exists := userStates[chatID]; !exists {
			userStates[chatID] = &UserState{}
		}

		user := userStates[chatID]

		switch {
		//Если пользователь отправил старт
		case messageText == "/start":
			handler.HandlerStart(bot, update.Message)

		// Если пользователь отправил /login
		case messageText == "/login":
			database.InitDB()
			if !database.IsExists(update.Message.From.UserName) {
				user.awaitingLogin = true
				msg := tgbotapi.NewMessage(chatID, "🔑Введите логин:")
				bot.Send(msg)
			} else {
				l, p := database.Select(update.Message.From.UserName)
				handler.HandlerLogin(bot, update.Message, l, p)
			}

		// Если бот ждёт логин
		case user.awaitingLogin:
			user.login = messageText
			user.awaitingLogin = false
			user.awaitingPassword = true
			msg := tgbotapi.NewMessage(chatID, "🔒Теперь введите пароль:")
			bot.Send(msg)

		// Если бот ждёт пароль → запускаем авторизацию
		case user.awaitingPassword:
			user.password = messageText
			user.awaitingPassword = false

			// Вызываем обработчик авторизации
			if handler.HandlerLogin(bot, update.Message, user.login, user.password) {
				database.Insert(update.Message.From.UserName, user.login, user.password)
			} else {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка авторизации, проверьте данные и нажмите /login")
				bot.Send(reply)
			}

		case messageText == "/update":
			database.InitDB()
			if database.IsExists(update.Message.From.UserName) {
				user.updateLogin = true
				msg := tgbotapi.NewMessage(chatID, "🔑Введите логин:")
				bot.Send(msg)
			}
		case user.updateLogin:
			user.login = messageText
			user.updateLogin = false
			user.uppdatePassword = true
			msg := tgbotapi.NewMessage(chatID, "🔒Теперь введите пароль:")
			bot.Send(msg)

		// Если бот ждёт пароль → запускаем авторизацию
		case user.uppdatePassword:
			user.password = messageText
			user.uppdatePassword = false

			// Вызываем обработчик авторизации
			if handler.HandlerLogin(bot, update.Message, user.login, user.password) {
				database.Update(update.Message.From.UserName, user.login, user.password)
			} else {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "❌Ошибка авторизации, проверьте данные и нажмите /login")
				bot.Send(reply)
			}
		// Любое другое сообщение
		default:
			msg := tgbotapi.NewMessage(chatID, "❗️Отправьте /login для БРС.")
			bot.Send(msg)
		}
	}
}
