package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type FolderMonitor struct {
	folder     string
	duration   time.Duration
	dispatcher Dispatcher
	ticker     *time.Ticker
}

func (m *FolderMonitor) Start() {
	m.ticker = time.NewTicker(m.duration)
	go m.run()
}

func (m *FolderMonitor) run() {
loop:
	for {
		select {
		case _, ok := <-m.ticker.C:
			if !ok {
				break loop
			}

			fis, err := ioutil.ReadDir(m.folder)
			if err != nil {
				//TODO: handle fatal error
				log.Println("FileMonitor: ", err)
			}

			if len(fis) == 0 {
				continue
			}

			for _, fi := range fis {
				filePath := filepath.Join(m.folder, fi.Name())
				log.Println("Found Mail:", filePath)
				m.dispatcher.Dispatch(filePath)
			}

			if waiter, ok := m.dispatcher.(Waiter); ok {
				waiter.Wait()
			}
		}
	}
}

func (m *FolderMonitor) Stop() {
	m.ticker.Stop()
}

func NewFolderMonitor(folder string, duration time.Duration, dispatcher Dispatcher) *FolderMonitor {
	return &FolderMonitor{folder: folder, duration: duration, dispatcher: dispatcher}
}
