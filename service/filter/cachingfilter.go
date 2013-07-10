package filter

import (
	"container/list"
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"sync"
)

type CachingFilter struct {
	next   Filter
	rwm    *sync.RWMutex
	l      *list.List
	size   int
	total  metrics.Counter
	hit    metrics.Counter
	cached metrics.Gauge
}

type Tuple struct {
	subject string
	result  Result
}

func (cf *CachingFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cf, mail.Name())
	cf.total.Inc(1)

	subject := mail.Subject()

	cf.rwm.RLock()
	for e := cf.l.Front(); e != nil; e = e.Next() {
		if t, ok := e.Value.(*Tuple); ok {
			if subject == t.subject {
				cf.rwm.RUnlock()
				cf.hit.Inc(1)
				return t.result
			}
		}
	}
	cf.rwm.RUnlock()

	result := cf.next.Filter(mail)

	cf.rwm.Lock()
	cf.l.PushFront(&Tuple{subject, result})
	if cf.l.Len() > cf.size {
		cf.l.Remove(cf.l.Back())
	}
	cf.cached.Update(int64(cf.l.Len()))
	cf.rwm.Unlock()

	return result
}

func (cf *CachingFilter) String() string {
	return "CachingFilter"
}

func NewCachingFilter(next Filter, size int) Filter {
	total := metrics.NewCounter()
	hit := metrics.NewCounter()
	cached := metrics.NewGauge()
	metrics.Register("CachingFilter-Total", total)
	metrics.Register("CachingFilter-Hit", hit)
	metrics.Register("CachingFilter-Cached", cached)
	return &CachingFilter{next, &sync.RWMutex{}, list.New(), size, total, hit, cached}
}
