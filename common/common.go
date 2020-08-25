package common

import (
	"encoding/json"
	"os"
	"tamastudy_news_crawler/domain/model"
	"unicode/utf8"
)

const (
	NewsCount      = 10
	LayoutYYYYMMDD = "2006-01-02"
	HttpsUrl       = `https://`
	Href           = "href"
)

func MinimizeContent(original string, limit int) string {
	ret := ""
	//Minimize String By limit Rune
	for _, r := range original {
		ret += string(r)

		if utf8.RuneCountInString(ret) >= limit {
			break
		}
	}

	ret += "..."

	return ret
}

func NewsSave(news []*model.News, fileName string) error{
	for _, n := range news{
		j, _ := json.Marshal(&n)
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			return err
		}
		if _, err := file.Write(j); err != nil {
			return err
		}

		if _, err := file.WriteString("\r\n\r\n"); err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}