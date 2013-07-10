package filter

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
)

type DeliverFilter struct {
	next  Filter
	total metrics.Counter
}

func (df *DeliverFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", df, mail.Name())
	df.total.Inc(1)

	result := df.next.Filter(mail)
	log.Printf("Move to %s, Mail:%s\n", result, mail.Name())
	err := common.MoveFile(mail.Path(), string(result))
	if err != nil {
		log.Printf("DeliverFilter: Err:%v, Mail:%s\n", err, mail.Name())
	}

	return Result("")
}

func (df *DeliverFilter) String() string {
	return "DeliverFilter"
}

func NewDeliverFilter(next Filter) Filter {
	total := metrics.NewCounter()
	metrics.Register("DeliverFilter-Total", total)
	return &DeliverFilter{next, total}
}
