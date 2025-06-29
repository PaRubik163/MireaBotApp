package lksMirea

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"resty.dev/v3"
	"strings"
)

func (p *Person) takeFIO(resp *resty.Response) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		logrus.Fatalf("Ошибка создания файла", err)
	}

	p.Name = (doc.Find(".ml-6").Find("h1").Text())
}
