package htmlutil

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
)

func ExtractText(reader io.Reader, remover func(string) (string, error)) (string, error) {
	z := html.NewTokenizer(reader)

	var buf bytes.Buffer
	bodyBlock := false

loop:
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.StartTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = true
			}
		case html.EndTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = false
			}
		case html.TextToken:
			if bodyBlock {
				buf.Write(z.Text())
			}
		case html.ErrorToken:
			if z.Err() != io.EOF {
				return "", z.Err()
			}
			break loop
		}
	}

	return remover(buf.String())
}
