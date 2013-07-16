package service

import (
	"github.com/parkghost/spamdefender/postfix"
	"time"
)

type Flusher interface {
	Flush() error
}

type PostfixFlusher struct {
	qmgrService string
	timeout     time.Duration
}

func (pf *PostfixFlusher) Flush() error {
	return postfix.MailTrigger(pf.qmgrService, postfix.QMGR_REQ_SCAN_INCOMING, pf.timeout)
}

func NewPostfixFlusher(service string, timeout time.Duration) Flusher {
	return &PostfixFlusher{service, timeout}
}
