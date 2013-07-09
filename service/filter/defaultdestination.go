package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type DefaultFilter struct {
	destFolder string
}

func (fdh *DefaultFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", fdh, mail.Name())

	return Result(fdh.destFolder + ps + mail.Name())
}

func (fdh *DefaultFilter) String() string {
	return "DefaultFilter"
}

func NewDestination(destFolder string) Filter {
	return &DefaultFilter{destFolder}
}
