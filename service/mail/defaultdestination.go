package mail

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type DefaultDestinationHandler struct {
	destFolder string
}

func (fdh *DefaultDestinationHandler) Handle(mail mailfile.Mail) bool {
	err := common.MoveFile(mail.Path(), fdh.destFolder+ps+mail.Name())
	if err != nil {
		log.Printf("DefaultDestinationHandler: Err: %v, Mail:%s\n", err, mail.Name())
	}

	return false
}

func (fdh *DefaultDestinationHandler) String() string {
	return "DefaultDestinationHandler"
}

func NewDefaultDestination(destFolder string) Handler {
	return &DefaultDestinationHandler{destFolder}
}
