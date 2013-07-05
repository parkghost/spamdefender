package mailhandler

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	"log"
	"os"
)

// THINK: a better chaining structure to avoid lost mail
// TODO: check duplicate mail

const ps = string(os.PathSeparator)

type MailHandler interface {
	Handle(mailfile.Mail) bool
}

type MailHandlerChain struct {
	handlers []MailHandler
}

func (mhl MailHandlerChain) Handle(mail mailfile.Mail) bool {
	for _, handler := range mhl.handlers {
		log.Printf("Run %s for Mail: %s\n", handler, mail.Name())
		passed := handler.Handle(mail)
		log.Printf("Passed: %t, Mail: %s\n", passed, mail.Name())

		if !passed {
			return true
		}
	}
	return true
}

func NewHandlerChain(list ...MailHandler) MailHandler {
	return &MailHandlerChain{list}
}

type MailHandlerAdapter struct {
	handler MailHandler
}

func (mha *MailHandlerAdapter) Handle(filePath string) {
	// THINK: implements factory method?
	mail := mailfile.NewPostfixMail(filePath)
	if err := mail.Parse(); err != nil {
		log.Printf("MailHandlerAdapter: Err: %v, Mail:%s\n", err, filePath)
		return
	}

	mha.handler.Handle(mail)
}

func NewMailHandlerAdapter(handler MailHandler) service.Handler {
	return &MailHandlerAdapter{handler}
}
