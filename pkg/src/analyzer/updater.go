package analyzer

import (
	"github.com/parkghost/bayesian"
	"log"
	"sync"
	"time"
)

type Updater interface {
	Update()
}

type DelayedUpdater struct {
	classifier          *bayesian.Classifier
	traningDataFilePath string
	delay               time.Duration

	rwm    *sync.RWMutex
	active chan bool
}

func (bu *DelayedUpdater) Update() {
	select {
	case bu.active <- true:
		go func() {
			time.Sleep(bu.delay)

			log.Printf("Updating classifier %s", bu.traningDataFilePath)
			bu.rwm.Lock()
			bu.classifier.WriteToFile(bu.traningDataFilePath)
			bu.rwm.Unlock()
			log.Printf("Updated classifier %s", bu.traningDataFilePath)

			<-bu.active
		}()
	default:
	}
}

func NewDelayedUpdater(classifier *bayesian.Classifier, traningDataFilePath string, updateDelay time.Duration, coordinator *sync.RWMutex) *DelayedUpdater {
	return &DelayedUpdater{classifier, traningDataFilePath, updateDelay, coordinator, make(chan bool, 1)}
}
