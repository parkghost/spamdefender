package mailfile

import (
	"bufio"
	"github.com/parkghost/pkg/net/mail"
	"io"
	"os"
	"path"
	"strings"
)

type plainTextMail struct {
	filePath string
	fd       *os.File
	subject  string
	content  io.Reader
	from     *mail.Address
	to       []*mail.Address
}

func (m *plainTextMail) Name() string {
	_, name := path.Split(m.filePath)
	return name
}

func (m *plainTextMail) Path() string {
	return m.filePath
}

func (m *plainTextMail) Subject() string {
	return m.subject
}

func (m *plainTextMail) Content() io.Reader {
	return m.content
}

func (m *plainTextMail) From() *mail.Address {
	return m.from
}

func (m *plainTextMail) To() []*mail.Address {
	return m.to
}

func (m *plainTextMail) Parse() (err error) {

	m.fd, err = os.Open(m.filePath)
	if err != nil {
		return
	}

	r := bufio.NewReader(m.fd)

	m.subject, err = r.ReadString('\n')
	if err != nil {
		return
	}
	m.subject = strings.TrimSpace(m.subject)

	fromText, err := r.ReadString('\n')
	if err != nil {
		return
	}
	m.from, err = mail.ParseAddress(strings.TrimSpace(fromText))
	if err != nil {
		return
	}

	toTextList, err := r.ReadString('\n')
	if err != nil {
		return
	}
	m.to, err = mail.ParseAddressList(strings.TrimSpace(toTextList))
	if err != nil {
		return
	}

	m.content = r
	return
}

func (m *plainTextMail) Close() error {
	if m.fd != nil {
		return m.fd.Close()
	}
	return nil
}

func (m *plainTextMail) String() string {
	return m.filePath
}

func NewPlainTextMail(filePath string) Mail {
	return &plainTextMail{filePath: filePath}
}
