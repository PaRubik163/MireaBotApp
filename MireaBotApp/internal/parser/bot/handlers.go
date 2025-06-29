package bot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	attend "mireabot/internal/parser/attendance"
	lk "mireabot/internal/parser/lksMirea"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	lks "mireabot/internal/parser/lksMirea"
)

func HandlerLogin(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, login, password string) bool {
	person := &lk.Person{}
	if !isGoodLogin(login) || !isGoodPassword(password) {
		return false
	}
	if !lks.Loginned(person, login, password) {
		return false
	} else {
		sometext := "Здравствуйте, " + person.Name + ", идёт авторизация..."

		reply := tgbotapi.NewMessage(msg.Chat.ID, sometext)
		sentMsg, err := bot.Send(reply)

		if err != nil {
			logrus.Fatalf("Ошибка отправки сообщения HandlerLogin", err)
		}

		time.Sleep(1 * time.Second)
		editGoodAutarizationMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, sentMsg.MessageID, "✅Авторизация успешна!")

		editGoodAutarization, err := bot.Send(editGoodAutarizationMsg)
		if err != nil {
			logrus.Fatalf("Ошибка редактирования сообщения об авторизации", err)
		}
		//Если мы успешно авторизировались в СДО, то логинимся на сайте посещений
		client := resty.New()

		err = attend.Logging(client, login, password)
		if err != nil {
			logrus.Fatal(err)
		}

		//авторизация->gRPC запрос на первый сервис, чтобы получить ID сервис->gRPC запрос, чтобы получить данные RatingScore
		//Конечно, нужно декодировать в структуру из прото, но я пока не понимаю как
		res, ok := attend.ParseGrpcResponse(client)
		if !ok {
			time.Sleep(1 * time.Second)
			editBadAutarizationMsg := tgbotapi.NewEditMessageText(msg.Chat.ID, sentMsg.MessageID, "❌Ошибка поиска предметов и баллов. Приносим свои извинения!")
			_, err = bot.Request(editBadAutarizationMsg)

			if err != nil {
				logrus.Fatalf("Ошибка удаления сообщения", err)
			}
		}

		//Здесь я все названия и сумму баллов по каждому предмету заношу в одно сообщение
		message := ""
		for _, item := range res {
			name, ok := item["name"].(string) //Названия предметов
			if !ok {
				logrus.Fatal("Нет такого поля name")
			}

			current_control, ok := item["current_control"].(float64) //Семестровый контроль
			if !ok {
				logrus.Fatal("Нет такого поля current_control")
			}

			attendance, ok := item["attendance"].(float64) //Баллы за посещаемость
			if !ok {
				logrus.Fatal("Нет такого поля attendance")
			}

			sum := current_control + attendance
			//Окрашивание
			switch {
			case sum >= 40:
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🔋")
			case sum < 40 && sum >= 25:
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🔶")
			case sum < 25:
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🚫")
			}
		}

		keyboard := buttonsForGoodAutarization()
		lastReply := tgbotapi.NewEditMessageText(msg.Chat.ID, editGoodAutarization.MessageID, message+"\n\n👉Логин: "+login+"\n🤐Пароль: "+password)
		lastReply.ReplyMarkup = &keyboard

		if _, err := bot.Send(lastReply); err != nil {
			logrus.Fatal("Ошибка отправки сообщения с баллами")
		}

		return true
	}
}

func isGoodLogin(login string) bool {
	if !strings.Contains(login, "@edu.mirea.ru") {
		return false
	}
	return true
}

func isGoodPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func BadAutorization(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "❌Неверные данные, проверьте логин и/или пароль!")
	reply.ReplyMarkup = buttonsForBadAutarization()

	if _, err := bot.Send(reply); err != nil {
		logrus.Fatalf("Ошибка отправки сообщения", err)
	}
}
