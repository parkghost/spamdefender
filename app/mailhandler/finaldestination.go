package mailhandler

import (
	"log"
	"spamdefender/common"
	"spamdefender/mailfile"
)

type FinalDestinationMailHandler struct {
	destFolder string
}

func (fdmh *FinalDestinationMailHandler) Handle(mail mailfile.Mail) bool {
	err := common.MoveFile(mail.Path(), fdmh.destFolder+ps+mail.Name())
	if err != nil {
		log.Printf("FinalDestinationMailHandler: Err: %v, Mail:%s\n", err, mail.Name())
	}
	return false
}

func (fdmh *FinalDestinationMailHandler) String() string {
	return "FinalDestinationMailHandler"
}

func NewFinalDestination(destFolder string) MailHandler {
	return &FinalDestinationMailHandler{destFolder}
}
