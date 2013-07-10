package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type RelayOnlyFilter struct {
	next        Filter
	localDomain string
	destFolder  string
}

func (sof *RelayOnlyFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", sof, mail.Name())

	sendOut := false
	for _, address := range mail.To() {
		if !strings.HasSuffix(address.Address, sof.localDomain) {
			sendOut = true
			break
		}
	}

	if !sendOut {
		return Result(sof.destFolder + ps + mail.Name())
	}

	return sof.next.Filter(mail)
}

func (sof *RelayOnlyFilter) String() string {
	return "RelayOnlyFilter"
}

func NewRelayOnlyFilter(next Filter, localDomain string, destFolder string) Filter {
	return &RelayOnlyFilter{next, localDomain, destFolder}
}
