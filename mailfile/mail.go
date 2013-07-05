package mailfile

import (
	"bytes"
	iconv "github.com/djimenez/iconv-go"
	"net/mail"
	"strings"
)

// TODO: write unit-test for To and From methods
// THINK: From and To support multi-charset

type Mail interface {
	Name() string
	Path() string
	Subject() string
	Content() string
	From() *mail.Address
	To() []*mail.Address
	Parse() error
}

func parseSubject(message *mail.Message) (string, error) {
	return DecodeRFC2047String(message.Header.Get("Subject"))
}

func parseFromAddress(message *mail.Message) (*mail.Address, error) {
	return mail.ParseAddress(message.Header.Get("From"))
}

func parseToAddress(message *mail.Message) ([]*mail.Address, error) {
	return mail.ParseAddressList(message.Header.Get("To"))
}

func parseBoby(message *mail.Message) (text string, err error) {
	//Content-Type: text/html;charset=UTF-8
	contentType := message.Header.Get("Content-Type")
	charset := contentType[strings.LastIndex(contentType, "=")+1:]

	reader := message.Body
	if charset != "UTF-8" {
		reader, err = iconv.NewReader(message.Body, charset, "UTF-8")
		if err != nil {
			return
		}
	}

	bodyContent := &bytes.Buffer{}
	_, err = bodyContent.ReadFrom(reader)
	if err != nil {
		return
	}

	text = bodyContent.String()
	return
}
