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

func (ch *CacheFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", ch, mail.Name())

	subject := mail.Subject()

	ch.rwm.RLock()
	for e := ch.l.Front(); e != nil; e = e.Next() {
		if t, ok := e.Value.(*Tuple); ok {
			if subject == t.subject {
				ch.rwm.RUnlock()
				return t.result
			}
		}
	}
	ch.rwm.RUnlock()

	result := ch.next.Filter(mail)

	ch.rwm.Lock()
	ch.l.PushFront(&Tuple{subject, result})
	if ch.l.Len() > ch.size {
		ch.l.Remove(ch.l.Back())
	}
	ch.rwm.Unlock()

	return result
}

func (ch *CacheFilter) String() string {
	return "CacheFilter"
}

func NewCache(next Filter, size int) Filter {
	return &CacheFilter{next, &sync.RWMutex{}, list.New(), size}
}
