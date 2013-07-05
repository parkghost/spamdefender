package app

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

const ps = string(os.PathSeparator)

type FolderMonitor struct {
	folder   string
	duration time.Duration
	handler  Handler
	ticker   *time.Ticker
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

			for _, fi := range fis {

				log.Println("Found Mail:", m.folder+ps+fi.Name())
				m.handler.Handle(m.folder + ps + fi.Name())
			}
		}
	}
}

func (m *FolderMonitor) Stop() {
	m.ticker.Stop()
}

func NewFolderMonitor(folder string, duration time.Duration, handler Handler) *FolderMonitor {
	return &FolderMonitor{folder: folder, duration: duration, handler: handler}
}