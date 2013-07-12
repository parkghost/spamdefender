package filter

import (
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	metrics "github.com/rcrowley/go-metrics"
	"log"
)

type ContentInspectionFilter struct {
	next     Filter
	allPass  bool
	anlz     analyzer.Analyzer
	total    metrics.Counter
	counters map[string]metrics.Counter
}

func (cif *ContentInspectionFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", cif, mail.Name())
	cif.total.Inc(1)

	htmlText := mail.Content()
	content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
	if err != nil {
		log.Printf("ContentInspectionFilter: Err:%v, Mail:%s\n", err, mail.Name())
		return cif.next.Filter(mail)
	}

	class := cif.anlz.Test(content)
	cif.counters[class].Inc(1)
	if cif.allPass || analyzer.Good == class {
		return cif.next.Filter(mail)
	}

	return Quarantine
}

func (cif *ContentInspectionFilter) String() string {
	return "ContentInspectionFilter"
}

func NewContentInspectionFilter(next Filter, allPass bool, traningDataFilePath string, dictDataFilePath string) Filter {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	total := metrics.NewCounter()
	metrics.Register("ContentInspectionFilter-Total", total)
	counters := make(map[string]metrics.Counter)
	for _, class := range []string{analyzer.Good, analyzer.Bad, analyzer.Neutral} {
		counter := metrics.NewCounter()
		counters[class] = counter
		metrics.Register("ContentInspectionFilter-"+class, counter)
	}
	return &ContentInspectionFilter{next, allPass, anlz, total, counters}
}
