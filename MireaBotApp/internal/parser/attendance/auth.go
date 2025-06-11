package attendance

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"net/http/cookiejar"
	"strings"
)

func Logging(client *resty.Client, login, password string) error {
	client.SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1; like Mac OS X) AppleWebKit/602.3 (KHTML, like Gecko)  Chrome/49.0.3440.106 Mobile Safari/601.9")
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

	jar, _ := cookiejar.New(nil)
	client.SetCookieJar(jar)

	if _, err := client.R().Get("https://attendance-app.mirea.ru/"); err != nil {
		return errors.New("Ошибка первого GET запроса на attendance")
	}

	//csrf токен
	resp, err := client.R().Get("https://attendance.mirea.ru/api/auth/login?redirectUri=https%3A%2F%2Fattendance-app.mirea.ru&rememberMe=True")

	if err != nil {
		return errors.New("Ошибка второго GET запроса на attendance")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		return errors.New("Ошибка создания файла на attendance")
	}

	csrfToken := doc.Find("input[name='csrfmiddlewaretoken']").AttrOr("value", "#")
	nextToken := doc.Find("input[name='next']").AttrOr("value", "#")

	_, err = client.R().
		SetFormData(map[string]string{
			"csrfmiddlewaretoken": csrfToken,
			"login":               login,
			"password":            password,
			"next":                nextToken,
		}).
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetHeader("origin", "https://login.mirea.ru").
		SetHeader("referer", "https://login.mirea.ru/login/?next=/oauth2/v1/authorize/%3Fclient_id%3DRkDSYWk7OPYsJ3KVehRbHRfjxdjIgmiCJ8j8IdvO8%26scope%3Dbasic%26response_type%3Dcode%26redirect_uri%3Dhttps%253A%252F%252Fattendance.mirea.ru%252Fapi%252Fmireaauth%26state%3Doauth_state%253A01970632-7fed-747a-a1c6-7dda25f047f1").
		Post("https://login.mirea.ru/login/")

	if err != nil {
		return errors.New("Ошибка POST запроса на attendance")
	}

	return nil
}
