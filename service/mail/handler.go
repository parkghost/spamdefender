package mail

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	"log"
	"os"
	"path"
)

// TODO: check duplicate mail

const ps = string(os.PathSeparator)

type Handler interface {
	Handle(mailfile.Mail) bool
}

type HandlerChain struct {
	handlers []Handler
}

func (hc HandlerChain) Handle(mail mailfile.Mail) bool {
	for _, handler := range hc.handlers {
		log.Printf("Run %s for Mail: %s\n", handler, mail.Name())
		next := handler.Handle(mail)
		if !next {
			break
		}
	}
	return true
}

func NewHandlerChain(list ...Handler) Handler {
	return &HandlerChain{list}
}

type FileHandlerAdapter struct {
	handler Handler
	factory mailfile.MailFileFactory
}

func (fha *FileHandlerAdapter) Handle(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		mail := fha.factory.Create(filePath)

		if err := mail.Parse(); err != nil {
			_, mailName := path.Split(filePath)
			log.Printf("FileHandlerAdapter: Err: %v, Mail:%s\n", err, mailName)
			return
		}

		fha.handler.Handle(mail)
	}
}

func NewFileHandlerAdapter(handler Handler, factory mailfile.MailFileFactory) service.Handler {
	return &FileHandlerAdapter{handler, factory}
}
