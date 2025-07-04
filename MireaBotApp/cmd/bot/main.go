package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/joho/godotenv"
	handler "mireabot/internal/parser/bot"
	"mireabot/internal/parser/bot/admin"
	database "mireabot/internal/parser/bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	login            string
	password         string
	awaitingLogin    bool
	awaitingPassword bool
	isUpdate         bool
}

var userStates = make(map[int64]*UserState)
var key []byte

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Ошибка загрузки .env файла", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	keyStr := os.Getenv("key")

	if botToken == "" || len(keyStr) < 32 {
		logrus.Fatalf("Short BotToken or key error!", err)
	}
	key = []byte(keyStr)
	bot, err := tgbotapi.NewBotAPI(botToken)

	//bot.Debug = true
	if err != nil {
		logrus.Fatal(err)
	}

	infoMsg := fmt.Sprintf("Бот %s запущен ", bot.Self.UserName)
	logrus.Info(infoMsg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		// Обработка сообщений от пользователя (текст)
		if update.Message != nil {
			if update.Message.ReplyToMessage != nil{
				admin.HandlForceReply(bot, update.Message)
			}
			
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			// Если пользователь новый — создаём для него состояние
			if _, exists := userStates[chatID]; !exists {
				userStates[chatID] = &UserState{}
			}
			user := userStates[chatID]

			switch {
			case user.awaitingLogin:
				user.login = text
				user.awaitingLogin = false
				user.awaitingPassword = true

				bot.Send(tgbotapi.NewMessage(chatID, "🔒 Теперь введите пароль:"))

			case user.awaitingPassword:
				user.password = text
				user.awaitingPassword = false

				//// Вызываем обработчик авторизации
				if handler.HandlerLogin(bot, update.Message, user.login, user.password) {
					if user.isUpdate {
						go database.Update(update.Message.From.UserName, user.login, user.password, key)
						go admin.HandlerAdminIfUpdate(bot, update.Message.From.UserName, user.login)
						user.isUpdate = false
					} else {
						go database.Insert(int(update.Message.Chat.ID), update.Message.From.UserName, user.login, user.password, key)
						go admin.HandlerAdminIfLogin(bot, update.Message.From.UserName, user.login, user.password)
					}
				} else {
					if user.isUpdate {
						handler.BadAutorization(bot, update.Message)
					} else {
						errButton := tgbotapi.NewInlineKeyboardButtonData("Попробовать ещё раз", "login")
						row := tgbotapi.NewInlineKeyboardRow(errButton)
						keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

						reply := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ОШИБКА АВТОРИЗАЦИИ\n\n🤔Невалидный логин или пароль\n🙏Пожалуйста, проверьте данные и попробуйте ещё раз")
						reply.ReplyMarkup = keyboard

						if _, err := bot.Send(reply); err != nil {
							logrus.Fatal("Ошибка отправки сообщения об ошибке авторизации", err)
						}
					}
				}

			case text == "/start":
				handler.SendStartButtons(bot, chatID)
			
			case text == "/broadcast":
				admin.HandlerBroadcast(bot, update)

			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Напиши /start или нажми кнопку"))
			}
		}

		// Обработка нажатий на кнопки
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID

			if _, exists := userStates[chatID]; !exists {
				userStates[chatID] = &UserState{}
			}
			user := userStates[chatID]

			switch callback.Data {
			case "login":
				go func() { //для параллельной обработки пользователей
					if r := recover(); r != nil {
						logrus.Warn("panic в момент login")
					}
					database.InitDB()
					if !database.IsExists(callback.From.UserName) {
						user.awaitingLogin = true
						bot.Send(tgbotapi.NewMessage(chatID, "🔑Введите логин МИРЭА:"))
					} else {
						l, p := database.SelectLoginandPassword(callback.From.UserName, key)
						go handler.HandlerLogin(bot, callback.Message, l, p)
						go admin.HandlerAdminIfLogin(bot, update.CallbackQuery.From.UserName, l, p)
					}
				}()
			case "update":
				go func() { //для параллельной обработки пользователей
					if r := recover(); r != nil {
						logrus.Warn("panic во время update")
					}
					database.InitDB()
					if database.IsExists(callback.From.UserName) {
						user.awaitingLogin = true
						user.isUpdate = true
						bot.Send(tgbotapi.NewMessage(chatID, "🔑Введите логин друга:"))
					}
				}()
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		}
	}
}
