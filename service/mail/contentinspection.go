package mail

import (
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type ContentInspectionHandler struct {
	allPass          bool
	quarantineFolder string
	anlz             analyzer.Analyzer
}

func (cih *ContentInspectionHandler) Handle(mail mailfile.Mail) bool {
	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionHandler: Err: %v, Mail:%s\n", err, mail.Name())
		return true
	}

	class := cih.anlz.Test(content)
	if cih.allPass || analyzer.Good == class {
		return true
	}

	err = common.MoveFile(mail.Path(), cih.quarantineFolder+ps+mail.Name())
	if err != nil {
		log.Printf("ContentInspectionHandler: Err: %v, Mail:%s\n", err, mail.Name())
	}

	return false
}

func (cih *ContentInspectionHandler) String() string {
	return "ContentInspectionHandler"
}

func NewContentInspection(allPass bool, quarantineFolder string, traningDataFilePath string, dictFilePath string) Handler {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return &ContentInspectionHandler{allPass, quarantineFolder, anlz}
}
