package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
)

type DefaultDestinationFilter struct {
	destFolder string
	total      metrics.Counter
}

func (ddf *DefaultDestinationFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", ddf, mail.Name())
	ddf.total.Inc(1)
	return Result(ddf.destFolder + ps + mail.Name())
}

func (ddf *DefaultDestinationFilter) String() string {
	return "DefaultDestinationFilter"
}

func NewDefaultDestinationFilter(destFolder string) Filter {
	total := metrics.NewCounter()
	metrics.Register("DefaultDestinationFilter-Total", total)
	return &DefaultDestinationFilter{destFolder, total}
}
