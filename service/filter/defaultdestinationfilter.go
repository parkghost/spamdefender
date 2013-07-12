package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
)

type DefaultDestinationFilter struct {
	total metrics.Counter
}

func (ddf *DefaultDestinationFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", ddf, mail.Name())
	ddf.total.Inc(1)
	return Incoming
}

func (ddf *DefaultDestinationFilter) String() string {
	return "DefaultDestinationFilter"
}

func NewDefaultDestinationFilter() Filter {
	total := metrics.NewCounter()
	metrics.Register("DefaultDestinationFilter-Total", total)
	return &DefaultDestinationFilter{total}
}
