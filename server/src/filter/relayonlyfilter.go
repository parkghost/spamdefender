package filter

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
	"strings"
)

type RelayOnlyFilter struct {
	next        Filter
	localDomain string
	total       metrics.Counter
	numOfRelay  metrics.Counter
}

func (sof *RelayOnlyFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", sof, mail.Name())
	sof.total.Inc(1)

	sendOut := false
	for _, address := range mail.To() {
		if !strings.HasSuffix(address.Address, sof.localDomain) {
			sendOut = true
			sof.numOfRelay.Inc(1)
			break
		}
	}

	if !sendOut {
		return Incoming
	}

	return sof.next.Filter(mail)
}

func (sof *RelayOnlyFilter) String() string {
	return "RelayOnlyFilter"
}

func NewRelayOnlyFilter(next Filter, localDomain string) *RelayOnlyFilter {
	total := metrics.NewCounter()
	numOfRelay := metrics.NewCounter()
	metrics.Register("RelayOnlyFilter-Total", total)
	metrics.Register("RelayOnlyFilter-Relay", numOfRelay)
	return &RelayOnlyFilter{next, localDomain, total, numOfRelay}
}
