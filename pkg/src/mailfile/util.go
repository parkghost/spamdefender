package mailfile

import (
	"bytes"
	"github.com/parkghost/pkg/net/mail"
	"strings"
)

func DecodeRFC2047String(s string) (string, error) {
	lines := strings.Split(s, " ")
	buf := bytes.NewBufferString("")
	for _, line := range lines {
		if strings.HasPrefix(line, "=?") {
			word, err := mail.DecodeRFC2047Word(line)
			if err != nil {
				return "", err
			}
			buf.WriteString(word)
		} else {
			buf.WriteString(" " + line)
		}
	}

	return buf.String(), nil
}
