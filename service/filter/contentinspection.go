package filter

import (
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
)

type ContentInspectionFilter struct {
	next             Filter
	allPass          bool
	quarantineFolder string
	anlz             analyzer.Analyzer
}

func (cih *ContentInspectionFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cih, mail.Name())

	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionFilter: Err:%v, Mail:%s\n", err, mail.Name())
		return cih.next.Filter(mail)
	}

	class := cih.anlz.Test(content)
	if cih.allPass || analyzer.Good == class {
		return cih.next.Filter(mail)
	}

	return Result(cih.quarantineFolder + ps + mail.Name())
}

func (cih *ContentInspectionFilter) String() string {
	return "ContentInspectionFilter"
}

func NewContentInspection(next Filter, allPass bool, quarantineFolder string, traningDataFilePath string, dictDataFilePath string) Filter {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return &ContentInspectionFilter{next, allPass, quarantineFolder, anlz}
}
