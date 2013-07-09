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

func (cif *ContentInspectionFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cif, mail.Name())

	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionFilter: Err:%v, Mail:%s\n", err, mail.Name())
		return cif.next.Filter(mail)
	}

	class := cif.anlz.Test(content)
	if cif.allPass || analyzer.Good == class {
		return cif.next.Filter(mail)
	}

	return Result(cif.quarantineFolder + ps + mail.Name())
}

func (cif *ContentInspectionFilter) String() string {
	return "ContentInspectionFilter"
}

func NewContentInspection(next Filter, allPass bool, quarantineFolder string, traningDataFilePath string, dictDataFilePath string) Filter {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	return &ContentInspectionFilter{next, allPass, quarantineFolder, anlz}
}
