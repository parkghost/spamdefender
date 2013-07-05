package mailhandler

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type MatchedSubjectMailHandler struct {
	subjectPrefix string
	destFolder    string
}

func (msmh *MatchedSubjectMailHandler) Handle(mail mailfile.Mail) bool {
	if strings.HasPrefix(mail.Subject(), msmh.subjectPrefix) {
		return true
	} else {
		err := common.MoveFile(mail.Path(), msmh.destFolder+ps+mail.Name())
		if err != nil {
			log.Printf("SendOutMailHandler: Err: %v, Mail:%s\n", err, mail.Name())
		}
		return false
	}

}

func (msmh *MatchedSubjectMailHandler) String() string {
	return "MatchedSubjectMailHandler"
}

func NewMatchedSubject(subjectPrefix string, destFolder string) MailHandler {
	return &MatchedSubjectMailHandler{subjectPrefix, destFolder}
}
