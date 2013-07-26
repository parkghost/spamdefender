package service

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"sync"
)

type Dispatcher interface {
	Dispatch(string)
}

type Waiter interface {
	Wait()
}

type FileHandler interface {
	Handle(string)
}

type PooledDispatcher struct {
	handler   FileHandler
	flusher   Flusher
	wg        *sync.WaitGroup
	semaphore chan bool
	meter     metrics.Meter
	timer     metrics.Timer
	active    metrics.Gauge
}

func (d *PooledDispatcher) Dispatch(filePath string) {
	d.wg.Add(1)
	d.active.Update(int64(len(d.semaphore)))
	d.semaphore <- true
	go func() {
		d.meter.Mark(1)
		d.timer.Time(func() {
			d.handler.Handle(filePath)
		})
		d.wg.Done()
		<-d.semaphore
		d.active.Update(int64(len(d.semaphore)))
	}()
}

func (d *PooledDispatcher) Wait() {
	d.wg.Wait()
	if d.flusher != nil {
		err := d.flusher.Flush()
		if err != nil {
			log.Printf("PooledDispatcher: flusher Err:%v", err)
		}
	}
}

func NewPooledDispatcher(handler FileHandler, flusher Flusher, size int) *PooledDispatcher {
	meter := metrics.NewMeter()
	timer := metrics.NewTimer()
	active := metrics.NewGauge()
	metrics.Register("PooledDispatcher-ProcessTime", timer)
	metrics.Register("PooledDispatcher-Mail", meter)
	metrics.Register("PooledDispatcher-Active", active)
	return &PooledDispatcher{handler, flusher, &sync.WaitGroup{}, make(chan bool, size), meter, timer, active}
}
