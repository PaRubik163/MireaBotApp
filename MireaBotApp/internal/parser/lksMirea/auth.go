package lksMirea

import (
	fake "github.com/EDDYCJY/fake-useragent"
	"github.com/sirupsen/logrus"
	"resty.dev/v3"
)

func Loginned(person *Person, login, password string) bool {
	client := resty.New()
	client.SetHeader("User-Agent", fake.Random())

	if _, err := client.R().Get("https://lk.mirea.ru/auth.php"); err != nil {
		logrus.Warn("Ошибка GET-запроса на сайте lk.mirea (GET-запрос)", err)
		return false
	}

	data := map[string]string{
		"AUTH_FORM":     "Y",
		"TYPE":          "AUTH",
		"USER_LOGIN":    login,
		"USER_PASSWORD": password,
		"USER_REMEMBER": "Y",
	}

	resp, err := client.R().SetFormData(data).Post("https://lk.mirea.ru/auth.php?login=yes")

	if err != nil || resp.StatusCode() != 200 {
		logrus.Warn("Ошибка авторизации на сайте lk.mirea (POST-запрос)", err)
		return false
	}
	person.takeFIO(resp)
	if person.Name == "" {
		return false
	}
	return true
}
