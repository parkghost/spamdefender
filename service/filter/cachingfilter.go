package filter

import (
	"container/list"
	"fmt"
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"sync"
)

type CachingProxy struct {
	target Filter
	rwm    *sync.RWMutex
	l      *list.List
	size   int
	total  metrics.Counter
	hit    metrics.Counter
	cached metrics.Gauge
}

type CacheKey interface {
	Key(mail mailfile.Mail) string
}

type CacheEntry struct {
	key    string
	result Result
}

func (cf *CachingProxy) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cf, mail.Name())
	cf.total.Inc(1)

	keyer, ok := cf.target.(CacheKey)

	if !ok {
		log.Fatalf("%s should implement CacheKey interface\n", cf.target)
		return None
	}

	key := keyer.Key(mail)

	cf.rwm.RLock()
	for e := cf.l.Front(); e != nil; e = e.Next() {
		if t, ok := e.Value.(*CacheEntry); ok {
			if key == t.key {
				cf.rwm.RUnlock()
				cf.hit.Inc(1)
				return t.result
			}
		}
	}
	cf.rwm.RUnlock()

	result := cf.target.Filter(mail)

	cf.rwm.Lock()
	cf.l.PushFront(&CacheEntry{key, result})
	if cf.l.Len() > cf.size {
		cf.l.Remove(cf.l.Back())
	}
	cf.cached.Update(int64(cf.l.Len()))
	cf.rwm.Unlock()

	return result
}

func (cf *CachingProxy) String() string {
	return fmt.Sprintf("CachingProxy(%s)", cf.target)
}

func NewCachingProxy(target Filter, size int) *CachingProxy {
	_, ok := target.(CacheKey)
	if !ok {
		log.Fatalf("%s should implement CacheKey interface\n", target)
		return nil
	}

	total := metrics.NewCounter()
	hit := metrics.NewCounter()
	cached := metrics.NewGauge()

	targetName := fmt.Sprintf("%s", target)

	metrics.Register("CachingProxy("+targetName+")-Total", total)
	metrics.Register("CachingProxy("+targetName+")-Hit", hit)
	metrics.Register("CachingProxy("+targetName+")-Cached", cached)
	return &CachingProxy{target, &sync.RWMutex{}, list.New(), size, total, hit, cached}
}
