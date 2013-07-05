package mailfile

import (
	"bytes"
	"net/mail"
)

// TODO: write unit-test for To and From methods

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

func parseBoby(message *mail.Message) (string, error) {
	bodyContent := &bytes.Buffer{}
	_, err := bodyContent.ReadFrom(message.Body)

	if err != nil {
		return "", err
	}
	return bodyContent.String(), nil
}
