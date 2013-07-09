package filter

import (
	"container/list"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"sync"
)

type CacheFilter struct {
	next Filter
	rwm  *sync.RWMutex
	l    *list.List
	size int
}

type Tuple struct {
	subject string
	result  Result
}

func (cf *CacheFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cf, mail.Name())

	subject := mail.Subject()

	cf.rwm.RLock()
	for e := cf.l.Front(); e != nil; e = e.Next() {
		if t, ok := e.Value.(*Tuple); ok {
			if subject == t.subject {
				cf.rwm.RUnlock()
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
	cf.rwm.Unlock()

	return result
}

func (cf *CacheFilter) String() string {
	return "CacheFilter"
}

func NewCache(next Filter, size int) Filter {
	return &CacheFilter{next, &sync.RWMutex{}, list.New(), size}
}
