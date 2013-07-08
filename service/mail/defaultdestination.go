package mail

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type DefaultHandler struct {
	destFolder string
}

func (fdh *DefaultHandler) Handle(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", fdh, mail.Name())

	return Result(fdh.destFolder + ps + mail.Name())
}

func (fdh *DefaultHandler) String() string {
	return "DefaultHandler"
}

func NewDestination(destFolder string) Handler {
	return &DefaultHandler{destFolder}
}
