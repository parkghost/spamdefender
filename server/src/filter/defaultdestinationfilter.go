package filter

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
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

func NewDefaultDestinationFilter() *DefaultDestinationFilter {
	total := metrics.NewCounter()
	metrics.Register("DefaultDestinationFilter-Total", total)
	return &DefaultDestinationFilter{total}
}
