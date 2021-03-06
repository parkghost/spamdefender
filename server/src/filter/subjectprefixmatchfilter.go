package filter

import (
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"mailfile"
	"strings"
)

type SubjectPrefixMatchFilter struct {
	next            Filter
	subjectPrefixes []string
	total           metrics.Counter
	matched         metrics.Counter
}

func (spmf *SubjectPrefixMatchFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", spmf, mail.Name())
	spmf.total.Inc(1)

	matched := false
	for _, subjectPrefix := range spmf.subjectPrefixes {
		if strings.HasPrefix(mail.Subject(), subjectPrefix) {
			matched = true
			spmf.matched.Inc(1)
			break
		}
	}

	if !matched {
		return Incoming
	}

	return spmf.next.Filter(mail)
}

func (spmf *SubjectPrefixMatchFilter) String() string {
	return "SubjectPrefixMatchFilter"
}

func NewSubjectPrefixMatchFilter(next Filter, subjectPrefixes []string) *SubjectPrefixMatchFilter {
	total := metrics.NewCounter()
	matched := metrics.NewCounter()
	metrics.Register("SubjectPrefixMatchFilter-Total", total)
	metrics.Register("SubjectPrefixMatchFilter-Matched", matched)
	return &SubjectPrefixMatchFilter{next, subjectPrefixes, total, matched}
}
