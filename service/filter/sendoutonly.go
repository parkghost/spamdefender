package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SendOutOnlyFilter struct {
	next        Filter
	localDomain string
	destFolder  string
}

func (soh *SendOutOnlyFilter) Filter(mail mailfile.Mail) Result {
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

	return soh.next.Filter(mail)
}

func (soh *SendOutOnlyFilter) String() string {
	return "SendOutOnlyFilter"
}

func NewSendOutOnly(next Filter, localDomain string, destFolder string) Filter {
	return &SendOutOnlyFilter{next, localDomain, destFolder}
}
