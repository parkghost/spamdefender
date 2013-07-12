package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"path"
)

const (
	None       = Result(0)
	Incoming   = Result(1)
	Quarantine = Result(2)
)

const ps = string(os.PathSeparator)

type Result int

type Filter interface {
	Filter(mailfile.Mail) Result
}

type FileHandlerAdapter struct {
	filter  Filter
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
			_, mailName := path.Split(filePath)
			log.Printf("FileHandlerAdapter: Err:%v, Mail:%s\n", err, mailName)
			return
		}
		defer mail.Close()

		fha.filter.Filter(mail)
	}
}

func NewFileHandlerAdapter(filter Filter, factory mailfile.MailFileFactory) service.Handler {
	total := metrics.NewCounter()
	meter := metrics.NewMeter()
	metrics.Register("FileHandlerAdapter-Total", total)
	metrics.Register("FileHandlerAdapter-Size", meter)
	return &FileHandlerAdapter{filter, factory, total, meter}
}
