package updater

import (
	"analyzer"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
	"mailpost"
	"os"
)

type ContentInspectionUpdater struct {
	anlz      analyzer.Analyzer
	class     string
	total     metrics.Counter
	malformed metrics.Counter
}

func (cih *ContentInspectionUpdater) Update(mail mailfile.Mail) {
	if leaner, ok := cih.anlz.(analyzer.Learner); ok {
		log.Printf("Run %s, Mail:%s\n", cih, mail.Name())
		cih.total.Inc(1)

		post, err := mailpost.Parse(mail)
		mail.Close()
		if err != nil {
			cih.malformed.Inc(1)
			log.Printf("ContentInspectionUpdater: Err: %v, Mail:%s\n", err, mail.Path())
			return
		}

		leaner.Learn(post.Subject, cih.class)
		leaner.Learn(post.Content, cih.class)

		err = os.Remove(mail.Path())
		if err != nil {
			log.Println(err)
		}
	}
}

func (cih *ContentInspectionUpdater) String() string {
	return "ContentInspectionUpdater"
}

func NewContentInspectionUpdater(anlz analyzer.Analyzer, class string) *ContentInspectionUpdater {
	total := metrics.NewCounter()
	malformed := metrics.NewCounter()
	metrics.Register("ContentInspectionUpdater-Total", total)
	metrics.Register("ContentInspectionUpdater-Malformed", malformed)

	return &ContentInspectionUpdater{anlz, class, total, malformed}
}
