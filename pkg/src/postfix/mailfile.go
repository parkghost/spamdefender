package postfix

import (
	"github.com/parkghost/pkg/net/mail"
	"io"
	"mailfile"
	"os"
	"path"
)

type PostfixMail struct {
	filePath string
	subject  string
	content  io.Reader
	from     *mail.Address
	to       []*mail.Address
}

func (m *PostfixMail) Name() string {
	_, name := path.Split(m.filePath)
	return name
}

func (m *PostfixMail) Path() string {
	return m.filePath
}

func (m *PostfixMail) Subject() string {
	return m.subject
}

func (m *PostfixMail) Content() io.Reader {
	return m.content
}

func (m *PostfixMail) From() *mail.Address {
	return m.from
}

func (m *PostfixMail) To() []*mail.Address {
	return m.to
}

func (m *PostfixMail) Parse() (err error) {
	f, err := os.Open(m.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var message *mail.Message
	message, err = mail.ReadMessage(NewRecordReader(f))
	if err != nil {
		return
	}

	m.subject, err = mailfile.ParseSubject(message)
	if err != nil {
		return
	}

	m.from, err = mailfile.ParseFromAddress(message)
	if err != nil {
		return
	}

	m.to, err = mailfile.ParseToAddress(message)
	if err != nil {
		return
	}

	m.content, err = mailfile.ParseBoby(message)
	if err != nil {
		return
	}

	return
}

func (m *PostfixMail) Close() error {
	return nil
}

func (m *PostfixMail) String() string {
	return m.filePath
}

func NewPostfixMail(filePath string) *PostfixMail {
	return &PostfixMail{filePath: filePath}
}

type PostfixMailFileFactory struct{}

func (p *PostfixMailFileFactory) Create(filePath string) mailfile.Mail {
	return NewPostfixMail(filePath)
}
