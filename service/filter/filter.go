package filter

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	"log"
	"os"
	"path"
)

const ps = string(os.PathSeparator)

type Result string

type Filter interface {
	Filter(mailfile.Mail) Result
}

type FileHandlerAdapter struct {
	filter  Filter
	factory mailfile.MailFileFactory
}

func (fha *FileHandlerAdapter) Handle(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		mail := fha.factory.Create(filePath)

		if err = mail.Parse(); err != nil {
			_, mailName := path.Split(filePath)
			log.Printf("FileHandlerAdapter: Err:%v, Mail:%s\n", err, mailName)
			return
		}

		result := fha.filter.Filter(mail)
		log.Printf("Move to %s, Mail:%s\n", result, mail.Name())
		err = common.MoveFile(mail.Path(), string(result))
		if err != nil {
			log.Printf("FileHandlerAdapter: Err:%v, Mail:%s\n", err, mail.Name())
		}

	}
}

func NewFileHandlerAdapter(filter Filter, factory mailfile.MailFileFactory) service.Handler {
	return &FileHandlerAdapter{filter, factory}
}
