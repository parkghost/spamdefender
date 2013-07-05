package html

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"errors"
	"fmt"
	"strings"
)

func BannerRemover(lineSeparator string, skipTop int, skipBottom int) func(string) (string, error) {
	return func(text string) (string, error) {

		lines := strings.Split(text, "\n")

		var pos []int
		for no, line := range lines {
			if strings.TrimRight(line, " ") == lineSeparator {
				pos = append(pos, no)
			}
		}

		if len(pos) == 0 {
			return text, errors.New(fmt.Sprintf("html: cannot detect lineSeparator: %s", lineSeparator))
		}

		// TODO: extract parameter
		if len(pos) != 3 {
			return text, errors.New("html: malformed mail content")
		}

		top, bottom := pos[skipTop], pos[len(pos)-skipBottom-1]
		return strings.Join(lines[top+1:bottom-1], "\n"), nil
	}
}

func ExtractText(htmlText string, remover func(string) (string, error)) (string, error) {
	z := html.NewTokenizer(strings.NewReader(htmlText))

	var buf bytes.Buffer
	bodyBlock := false

loop:
	for {
		tokenType := z.Next()
		switch {
		case tokenType == html.StartTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = true
			}
		case tokenType == html.EndTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = false
			}
		case tokenType == html.TextToken:
			if bodyBlock {
				buf.Write(z.Text())
			}
		case tokenType == html.ErrorToken:
			break loop
		case z.Err() != nil:
			fmt.Println(z.Err())
			break loop
		}
	}

	return remover(buf.String())
}
