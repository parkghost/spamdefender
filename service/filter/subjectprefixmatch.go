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

func (msh *SubjectPrefixMatchFilter) Filter(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", msh, mail.Name())

	if !strings.HasPrefix(mail.Subject(), msh.subjectPrefix) {
		return Result(msh.destFolder + ps + mail.Name())
	}

	return msh.next.Filter(mail)
}

func (msh *SubjectPrefixMatchFilter) String() string {
	return "SubjectPrefixMatchFilter"
}

func NewSubjectPrefixMatch(next Filter, subjectPrefix string, destFolder string) Filter {
	return &SubjectPrefixMatchFilter{next, subjectPrefix, destFolder}
}
