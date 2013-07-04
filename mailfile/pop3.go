package mailfile

import (
	"bytes"
	"io"
	"net/mail"
	"os"
)

type POP3Mail struct {
	parsed   bool
	filePath string
	subject  string
	content  string
}

func (m *POP3Mail) Subject() string {
	if !m.parsed {
		m.parse()
	}

	return m.subject
}

func (m *POP3Mail) Content() string {
	if !m.parsed {
		m.parse()
	}

	return m.content
}

func (m *POP3Mail) parse() (err error) {
	fs, err := os.Open(m.filePath)
	if err != nil {
		return
	}
	defer fs.Close()

	message, err := mail.ReadMessage(fs)
	if err != nil {
		return
	}

	rawSubjectStr := message.Header.Get("Subject")

	m.subject, err = DecodeRFC2047String(rawSubjectStr)
	if err != nil {
		return
	}

	bodyContent := &bytes.Buffer{}
	io.Copy(bodyContent, message.Body)
	m.content = bodyContent.String()

	m.parsed = true

	return
}

func NewPOP3Mail(filePath string) Mail {
	return &POP3Mail{filePath: filePath}
}
