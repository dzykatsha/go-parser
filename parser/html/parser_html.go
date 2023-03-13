package html

import (
	"golang.org/x/net/html"
	"io/ioutil"
	"mvdan.cc/xurls"
	"strings"
)

func ReadHtmlFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseHtml(text string) []string {
	tokenizer := html.NewTokenizer(strings.NewReader(text))
	var (
		content []string
		links []string
		isList  bool
		isTable bool
	)
	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			if len(content) > 0 {
				content = removeEmptyStrings(content)
			}
			if len(links) > 0 {
				content = append(content, "\n")
				links = removeEmptyStrings(links)
				content = append(content, links...)
			}
			return content

		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			isList = (token.Data == "li")
			isTable = (token.Data == "td")
			if token.Data == "a" {
				links = append(links, xurls.Relaxed.FindString(token.String()))
			}

		case tokenType == html.TextToken:
			token := tokenizer.Token()
			if isList {
				content = append(content, token.Data)

			} else if isTable {
				text := (string)(token.Data)
				trimmed := strings.TrimSpace(text)
				content = append(content, trimmed)
			} else if strings.TrimSpace(token.Data) != "" {
				data := strings.Join(strings.Fields(token.Data), " ")
				content = append(content, data)
			}
			isList = false
			isTable = false
		}
	}

}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "\n" {
			r = append(r, str)
		}
	}
	return r
}
