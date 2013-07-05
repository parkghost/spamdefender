package mailfile

import (
	"net/mail"
	"os"
	"path"
)

type POP3Mail struct {
	filePath string
	subject  string
	content  string
	from     *mail.Address
	to       []*mail.Address
}

func (m *POP3Mail) Name() string {
	_, name := path.Split(m.filePath)
	return name
}

func (m *POP3Mail) Path() string {
	return m.filePath
}

func (m *POP3Mail) Subject() string {
	return m.subject
}

func (m *POP3Mail) Content() string {
	return m.content
}

func (m *POP3Mail) From() *mail.Address {
	return m.from
}

func (m *POP3Mail) To() []*mail.Address {
	return m.to
}

func (m *POP3Mail) Parse() (err error) {

	fs, err := os.Open(m.filePath)
	if err != nil {
		return
	}
	defer fs.Close()

	message, err := mail.ReadMessage(fs)
	if err != nil {
		return
	}

	m.subject, err = parseSubject(message)
	if err != nil {
		return
	}

	m.from, err = parseFromAddress(message)
	if err != nil {
		return
	}

	m.to, err = parseToAddress(message)
	if err != nil {
		return
	}

	m.content, err = parseBoby(message)
	if err != nil {
		return
	}

	return
}

func (m *POP3Mail) String() string {
	return m.filePath
}

func NewPOP3Mail(filePath string) Mail {
	return &POP3Mail{filePath: filePath}
}
