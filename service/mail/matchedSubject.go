package mail

import (
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/mailfile"
	"log"
	"strings"
)

type MatchedSubjectHandler struct {
	subjectPrefix string
	destFolder    string
}

func (msh *MatchedSubjectHandler) Handle(mail mailfile.Mail) bool {
	if !strings.HasPrefix(mail.Subject(), msh.subjectPrefix) {
		err := common.MoveFile(mail.Path(), msh.destFolder+ps+mail.Name())
		if err != nil {
			log.Printf("MatchedSubjectHandler: Err: %v, Mail:%s\n", err, mail.Name())
		}
		return false
	}

	return true
}

func (msh *MatchedSubjectHandler) String() string {
	return "MatchedSubjectHandler"
}

func NewMatchedSubject(subjectPrefix string, destFolder string) Handler {
	return &MatchedSubjectHandler{subjectPrefix, destFolder}
}
