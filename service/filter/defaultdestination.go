package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type DefaultDestinationFilter struct {
	destFolder string
}

func (ddf *DefaultDestinationFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", ddf, mail.Name())

	return Result(ddf.destFolder + ps + mail.Name())
}

func (ddf *DefaultDestinationFilter) String() string {
	return "DefaultDestinationFilter"
}

func NewDefaultDestination(destFolder string) Filter {
	return &DefaultDestinationFilter{destFolder}
}
