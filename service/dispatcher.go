package service

import (
	"sync"
)

type Dispatcher interface {
	Dispatch(string)
}

type Waiter interface {
	Wait()
}

type Handler interface {
	Handle(string)
}

type PooledDispatcher struct {
	handler   Handler
	wg        *sync.WaitGroup
	semaphore chan bool
}

func (d *PooledDispatcher) Dispatch(filePath string) {
	d.wg.Add(1)
	d.semaphore <- true
	go func() {
		d.handler.Handle(filePath)
		d.wg.Done()
		<-d.semaphore
	}()
}

func (d *PooledDispatcher) Wait() {
	d.wg.Wait()
}

func NewPooledDispatcher(handler Handler, size int) Dispatcher {
	return &PooledDispatcher{handler, &sync.WaitGroup{}, make(chan bool, size)}
}
