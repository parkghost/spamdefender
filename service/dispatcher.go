package service

import (
	"sync"
)

type Dispatcher struct {
	active    map[string]bool
	rwm       *sync.RWMutex
	handler   Handler
	semaphore chan bool
}

func (d *Dispatcher) Handle(filePath string) {
	d.semaphore <- true
	d.rwm.RLock()
	_, found := d.active[filePath]
	d.rwm.RUnlock()
	if !found {
		go func() {
			d.rwm.Lock()
			d.active[filePath] = true
			d.rwm.Unlock()

			d.handler.Handle(filePath)

			d.rwm.Lock()
			delete(d.active, filePath)
			d.rwm.Unlock()
			<-d.semaphore
		}()
	}
}

func NewDispatcher(handler Handler, size int) *Dispatcher {
	return &Dispatcher{make(map[string]bool), &sync.RWMutex{}, handler, make(chan bool, size)}
}
