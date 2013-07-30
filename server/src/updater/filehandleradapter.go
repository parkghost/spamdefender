package updater

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
	"os"
	"path/filepath"
)

type FileHandlerAdapter struct {
	updater Updater
	factory mailfile.MailFileFactory
	total   metrics.Counter
	meter   metrics.Meter
}

func (fha *FileHandlerAdapter) Handle(filePath string) {
	fha.total.Inc(1)
	f, err := os.Stat(filePath)
	if err == nil {
		fha.meter.Mark(f.Size())
		mail := fha.factory.Create(filePath)

		if err = mail.Parse(); err != nil {
			_, mailName := filepath.Split(filePath)
			log.Printf("FileHandlerAdapter: Err:%v, Mail:%s\n", err, mailName)
			return
		}
		defer mail.Close()

		fha.updater.Update(mail)
	}
}

func NewFileHandlerAdapter(updater Updater, factory mailfile.MailFileFactory) *FileHandlerAdapter {
	total := metrics.NewCounter()
	meter := metrics.NewMeter()
	metrics.Register("FileHandlerAdapter-Total", total)
	metrics.Register("FileHandlerAdapter-Size", meter)
	return &FileHandlerAdapter{updater, factory, total, meter}
}
