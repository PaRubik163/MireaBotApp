package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func SendStartButtons(bot *tgbotapi.BotAPI, chatID int64) {
	login := tgbotapi.NewInlineKeyboardButtonData("Авторизоваться", "login")

	row := tgbotapi.NewInlineKeyboardRow(login)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg := tgbotapi.NewMessage(
		chatID,
		"👋*Рады тебя видеть в __MireaScore__*\\!\n\n"+
			"📌Что делает этот бот\\?\n"+
			"🤖Этот бот авторизирует тебя на сайте МИРЭА\n\n"+
			"🔢Присылает баллы по каждой дисциплине\n\n"+
			"🤝P\\.S Исключительно для просмотра успеваемости\\!",
	)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		logrus.Fatalf("Ошибка отправки /start сообщения")
	}
}

func buttonsForGoodAutarization() tgbotapi.InlineKeyboardMarkup {
	oneMore, checkAnother := tgbotapi.NewInlineKeyboardButtonData("Мои баллы", "login"), tgbotapi.NewInlineKeyboardButtonData("Баллы друга", "update")

	rows := tgbotapi.NewInlineKeyboardRow(oneMore, checkAnother)

	return tgbotapi.NewInlineKeyboardMarkup(rows)
}

func buttonsForBadAutarization() tgbotapi.InlineKeyboardMarkup {
	checkMyScore, returnOneMore := tgbotapi.NewInlineKeyboardButtonData("Посмотреть свои баллы", "login"), tgbotapi.NewInlineKeyboardButtonData("Попробовать ещё раз", "update")

	rows := tgbotapi.NewInlineKeyboardRow(checkMyScore, returnOneMore)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows)

	return keyboard
}
