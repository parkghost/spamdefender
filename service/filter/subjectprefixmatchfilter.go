package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SubjectPrefixMatchFilter struct {
	next            Filter
	subjectPrefixes []string
	destFolder      string
}

func (spmf *SubjectPrefixMatchFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", spmf, mail.Name())

	matched := false
	for _, subjectPrefix := range spmf.subjectPrefixes {
		if strings.HasPrefix(mail.Subject(), subjectPrefix) {
			matched = true
			break
		}
	}

	if !matched {
		return Result(spmf.destFolder + ps + mail.Name())
	}

	return spmf.next.Filter(mail)
}

func (spmf *SubjectPrefixMatchFilter) String() string {
	return "SubjectPrefixMatchFilter"
}

func NewSubjectPrefixMatchFilter(next Filter, subjectPrefixes []string, destFolder string) Filter {
	return &SubjectPrefixMatchFilter{next, subjectPrefixes, destFolder}
}
