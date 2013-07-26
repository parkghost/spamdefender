package main

import (
	"postfix"
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

func NewPostfixFlusher(service string, timeout time.Duration) *PostfixFlusher {
	return &PostfixFlusher{service, timeout}
}
