package mail

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SendOutOnlyHandler struct {
	next        Handler
	localDomain string
	destFolder  string
}

func (soh *SendOutOnlyHandler) Handle(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", soh, mail.Name())

	sendOut := false
	for _, address := range mail.To() {
		if !strings.HasSuffix(address.Address, soh.localDomain) {
			sendOut = true
			break
		}
	}

	if !sendOut {
		return Result(soh.destFolder + ps + mail.Name())
	}

	return soh.next.Handle(mail)
}

func (soh *SendOutOnlyHandler) String() string {
	return "SendOutOnlyHandler"
}

func NewSendOutOnly(next Handler, localDomain string, destFolder string) Handler {
	return &SendOutOnlyHandler{next, localDomain, destFolder}
}
