package mail

import (
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type ContentInspectionHandler struct {
	next             Handler
	allPass          bool
	quarantineFolder string
	anlz             analyzer.Analyzer
}

func (cih *ContentInspectionHandler) Handle(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cih, mail.Name())

	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionHandler: Err:%v, Mail:%s\n", err, mail.Name())
		return cih.next.Handle(mail)
	}

	class := cih.anlz.Test(content)
	if cih.allPass || analyzer.Good == class {
		return cih.next.Handle(mail)
	}

	return Result(cih.quarantineFolder + ps + mail.Name())
}

func (cih *ContentInspectionHandler) String() string {
	return "ContentInspectionHandler"
}

func NewContentInspection(next Handler, allPass bool, quarantineFolder string, traningDataFilePath string, dictDataFilePath string) Handler {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return &ContentInspectionHandler{next, allPass, quarantineFolder, anlz}
}
