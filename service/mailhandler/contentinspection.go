package mailhandler

import (
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type ContentInspectionMailHandler struct {
	allPass          bool
	quarantineFolder string
	anlz             *analyzer.Analyzer
}

func (cimh *ContentInspectionMailHandler) Handle(mail mailfile.Mail) bool {
	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionMailHandler: Err: %v, Mail:%s\n", err, mail.Name())
		return false
	}

	score, pass := cimh.anlz.Test(content)

	// TODO: print readable score format
	log.Printf("Score: %v, Mail:%s\n", score, mail.Name())

	finalResult := false

	if cimh.allPass {
		finalResult = cimh.allPass
	} else {
		finalResult = pass
	}

	if !finalResult {
		common.MoveFile(mail.Path(), cimh.quarantineFolder+ps+mail.Name())
		return false
	}
	return finalResult
}

func (cimh *ContentInspectionMailHandler) String() string {
	return "ContentInspectionMailHandler"
}

func NewContentInspection(allPass bool, quarantineFolder string, traningDataFilePath string, dictFilePath string) MailHandler {
	anlz, err := analyzer.NewAnalyzer(traningDataFilePath, dictFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return &ContentInspectionMailHandler{allPass, quarantineFolder, anlz}
}
