package common

import (
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
