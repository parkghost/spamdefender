package mail

import (
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type SubjectPrefixMatchHandler struct {
	next          Handler
	subjectPrefix string
	destFolder    string
}

func (msh *SubjectPrefixMatchHandler) Handle(mail mailfile.Mail) Result {
	log.Printf("Run %s, Mail:%s\n", msh, mail.Name())

	if !strings.HasPrefix(mail.Subject(), msh.subjectPrefix) {
		return Result(msh.destFolder + ps + mail.Name())
	}

	return msh.next.Handle(mail)
}

func (msh *SubjectPrefixMatchHandler) String() string {
	return "SubjectPrefixMatchHandler"
}

func NewSubjectPrefixMatch(next Handler, subjectPrefix string, destFolder string) Handler {
	return &SubjectPrefixMatchHandler{next, subjectPrefix, destFolder}
}
