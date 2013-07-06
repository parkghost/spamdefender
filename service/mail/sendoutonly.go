package mail

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SendOutOnlyHandler struct {
	localDomain string
	destFolder  string
}

func (soh *SendOutOnlyHandler) Handle(mail mailfile.Mail) bool {
	sendOut := false
	for _, address := range mail.To() {
		if !strings.HasSuffix(address.Address, soh.localDomain) {
			sendOut = true
			break
		}
	}

	if !sendOut {
		err := common.MoveFile(mail.Path(), soh.destFolder+ps+mail.Name())
		if err != nil {
			log.Printf("SendOutOnlyHandler: Err: %v, Mail:%s\n", err, mail.Name())
		}
		return false
	}

	return true
}

func (soh *SendOutOnlyHandler) String() string {
	return "SendOutOnlyHandler"
}

func NewSendOutOnly(localDomain string, destFolder string) Handler {
	return &SendOutOnlyHandler{localDomain, destFolder}
}
