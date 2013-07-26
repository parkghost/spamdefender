package mailfile

import (
	"github.com/parkghost/pkg/net/mail"
	"io"
	"os"
	"path"
)

type POP3Mail struct {
	filePath string
	fd       *os.File
	subject  string
	content  io.Reader
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

func (m *POP3Mail) Content() io.Reader {
	return m.content
}

func (m *POP3Mail) From() *mail.Address {
	return m.from
}

func (m *POP3Mail) To() []*mail.Address {
	return m.to
}

func (m *POP3Mail) Parse() (err error) {

	m.fd, err = os.Open(m.filePath)
	if err != nil {
		return
	}

	message, err := mail.ReadMessage(m.fd)
	if err != nil {
		return
	}

	m.subject, err = ParseSubject(message)
	if err != nil {
		return
	}

	m.from, err = ParseFromAddress(message)
	if err != nil {
		return
	}

	m.to, err = ParseToAddress(message)
	if err != nil {
		return
	}

	m.content, err = ParseBoby(message)
	if err != nil {
		return
	}

	return
}

func (m *POP3Mail) Close() error {
	if m.fd != nil {
		return m.fd.Close()
	}
	return nil
}

func (m *POP3Mail) String() string {
	return m.filePath
}

func NewPOP3Mail(filePath string) *POP3Mail {
	return &POP3Mail{filePath: filePath}
}

type POP3MailFileFactory struct{}

func (p *POP3MailFileFactory) Create(filePath string) Mail {
	return NewPOP3Mail(filePath)
}
