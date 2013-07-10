package filter

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SubjectPrefixMatchFilter struct {
	next          Filter
	subjectPrefix string
	destFolder    string
}

func (spmf *SubjectPrefixMatchFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", spmf, mail.Name())

	if !strings.HasPrefix(mail.Subject(), spmf.subjectPrefix) {
		return Result(spmf.destFolder + ps + mail.Name())
	}

	return spmf.next.Filter(mail)
}

func (spmf *SubjectPrefixMatchFilter) String() string {
	return "SubjectPrefixMatchFilter"
}

func NewSubjectPrefixMatchFilter(next Filter, subjectPrefix string, destFolder string) Filter {
	return &SubjectPrefixMatchFilter{next, subjectPrefix, destFolder}
}
