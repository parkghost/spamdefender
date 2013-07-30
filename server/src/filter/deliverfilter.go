package filter

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
	"os"
	"path/filepath"
)

type DeliverFilter struct {
	next     Filter
	paths    map[Result]string
	total    metrics.Counter
	counters map[Result]metrics.Counter
}

func (df *DeliverFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", df, mail.Name())
	df.total.Inc(1)

	result := df.next.Filter(mail)

	df.counters[result].Inc(1)

	if result == None {
		log.Fatalf("DeliverFilter: the filter result should not be None, Mail:%s\n", mail.Name())
	}

	destination := filepath.Join(df.paths[result], mail.Name())
	log.Printf("Move to %s, Mail:%s\n", destination, mail.Name())

	err := os.Rename(mail.Path(), destination)
	if err != nil {
		log.Printf("DeliverFilter: Err:%v, Mail:%s\n", err, mail.Name())
	}

	return None
}

func (df *DeliverFilter) String() string {
	return "DeliverFilter"
}

func NewDeliverFilter(next Filter, paths map[Result]string) *DeliverFilter {
	total := metrics.NewCounter()
	metrics.Register("DeliverFilter-Total", total)

	counters := make(map[Result]metrics.Counter)
	for result, path := range paths {
		counter := metrics.NewCounter()
		counters[result] = counter
		metrics.Register("DeliverFilter-"+path, counter)
	}
	return &DeliverFilter{next, paths, total, counters}
}
