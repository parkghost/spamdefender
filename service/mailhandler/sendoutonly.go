package mailhandler

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SendOutOnlyMailHandler struct {
	localDomain string
	destFolder  string
}

func (somh *SendOutOnlyMailHandler) Handle(mail mailfile.Mail) bool {
	sendOut := false
	for _, address := range mail.To() {
		if !strings.HasSuffix(address.Address, somh.localDomain) {
			sendOut = true
			break
		}
	}

	if !sendOut {
		err := common.MoveFile(mail.Path(), somh.destFolder+ps+mail.Name())
		if err != nil {
			log.Printf("SendOutOnlyMailHandler: Err: %v, Mail:%s\n", err, mail.Name())
		}
		return false
	}

	return true
}

func (somh *SendOutOnlyMailHandler) String() string {
	return "SendOutOnlyMailHandler"
}

func NewSendOutOnly(localDomain string, destFolder string) MailHandler {
	return &SendOutOnlyMailHandler{localDomain, destFolder}
}
