package mailfile

import (
	"github.com/parkghost/pkg/net/mail"
	"io"
)

type Mail interface {
	Name() string
	Path() string
	Subject() string
	Content() io.Reader
	From() *mail.Address
	To() []*mail.Address
	Parse() error
	Close() error
}

func ParseSubject(message *mail.Message) (string, error) {
	return DecodeRFC2047String(message.Header.Get("Subject"))
}

func ParseFromAddress(message *mail.Message) (*mail.Address, error) {
	return mail.ParseAddress(message.Header.Get("From"))
}

func ParseToAddress(message *mail.Message) ([]*mail.Address, error) {
	return mail.ParseAddressList(message.Header.Get("To"))
}

func ParseBoby(message *mail.Message) (reader io.Reader, err error) {
	// THINK: handle variant Content-Type, Content-Transfer-Encoding, charset
	return message.Body, nil
}

type MailFileFactory interface {
	Create(string) Mail
}

type POP3MailFileFactory struct{}

func (p *POP3MailFileFactory) Create(filePath string) Mail {
	return NewPOP3Mail(filePath)
}
