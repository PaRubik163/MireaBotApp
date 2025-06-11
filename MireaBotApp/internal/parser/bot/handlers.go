package bot

// Надо придумать как рзабить эту функцию
import (
	"fmt"
	"log"
	attend "mireabot/internal/parser/attendance"
	lk "mireabot/internal/parser/lksMirea"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	lks "mireabot/internal/parser/lksMirea"
)

func HandlerStart(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	welcomeText := fmt.Sprintf("👋Рады тебя видеть в MireaScore!\n\n📌Что делает этот бот?\n🤖Этот бот авторизирует тебя на сайте МИРЭА\n\n🔢Присылает баллы по каждой дисциплине\nВведи команду /login для работы ЛКС\n\n🤝P.S Исключительно для просмотра успеваемости!")

	reply := tgbotapi.NewMessage(msg.Chat.ID, welcomeText)
	_, err := bot.Send(reply)

	if err != nil {
		log.Fatalf("Ошибка отправки сообщения HandlerStart", err)
	}
}

func HandlerLogin(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, login, password string) bool {
	person := &lk.Person{}
	if !IsGoodLogin(login) || !IsGoodPassword(password) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "❌Невалидный логин или пароль\n🙏Пожалуйста, проверьте данные и нажмите /update")
		bot.Send(reply)
		return false
	}
	if !lks.Loginned(person, login, password) {
		return false
	} else {
		sometext := "Здравствуйте, " + person.Name + ", идёт авторизация..."

		reply := tgbotapi.NewMessage(msg.Chat.ID, sometext)
		sentMsg, err := bot.Send(reply)

		if err != nil {
			log.Fatalf("Ошибка отправки сообщения HandlerLogin", err)
		}

		time.Sleep(2 * time.Second)
		deletemsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, sentMsg.MessageID)
		_, err = bot.Request(deletemsg)

		if err != nil {
			log.Fatalf("Ошибка удаления сообщения", err)
		}

		sentMsg, err = bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "✅Авторизация успешна!"))

		if err != nil {
			log.Fatalf("Ошибка отправки сообщения HandlerLogin", err)
		}
		//Если мы успешно авторизировались в СДО, то логинимся на сайте посещений
		client := resty.New()

		err = attend.Logging(client, login, password)
		if err != nil {
			log.Fatal(err)
		}

		//авторизация->gRPC запрос на первый сервис, чтобы получить ID сервис->gRPC запрос, чтобы получить данные RatingScore
		//Конечно, нужно декодировать в структуру из прото, но я пока не понимаю как
		res, ok := attend.ParseGrpcResponse(client)
		if !ok {
			time.Sleep(2 * time.Second)
			deletemsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, sentMsg.MessageID)
			_, err = bot.Request(deletemsg)

			if err != nil {
				log.Fatalf("Ошибка удаления сообщения", err)
			}

			reply := tgbotapi.NewMessage(msg.Chat.ID, "❌Ошибка поиска предметов и баллов. Приносим свои извинения!")
			bot.Send(reply)
		}

		//Здесь я все названия и сумму баллов по каждому предмету заношу в одно сообщение
		message := ""
		for _, item := range res {
			name, ok := item["name"].(string) //Названия предметов
			if !ok {
				log.Fatal("Нет такого поля name")
			}

			current_control, ok := item["current_control"].(float64) //Семестровый контроль
			if !ok {
				log.Fatal("Нет такого поля current_control")
			}

			attendance, ok := item["attendance"].(float64) //Баллы за посещаемость
			if !ok {
				log.Fatal("Нет такого поля attendance")
			}

			sum := current_control + attendance
			//Окрашивание
			if sum >= 40 {
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🟢")
			}
			if sum < 40 && sum >= 25 {
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🟡")
			}
			if sum < 25 {
				message += fmt.Sprintf("%s %.1f %s\n", name, sum, "🔴")
			}
		}

		lastReply := tgbotapi.NewMessage(msg.Chat.ID, message+"\n\n👉Просто нажмите /login для просмотра баллов")
		bot.Send(lastReply)

		return true
	}
}
