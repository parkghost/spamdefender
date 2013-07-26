package postfix

import (
	"net"
	"time"
)

// postfix-2.10.1/src/global/mail_proto.h
const (
	QMGR_REQ_SCAN_DEFERRED = 'D' // scan deferred queue
	QMGR_REQ_SCAN_INCOMING = 'I' // scan incoming queue
	QMGR_REQ_FLUSH_DEAD    = 'F' // flush dead xport/site
	QMGR_REQ_SCAN_ALL      = 'A' // ignore time stamps
)

// postfix-2.10.1/src/global/mail_trigger.c
// postfix-2.10.1/src/util/unix_trigger.c
func MailTrigger(service string, request rune, timeout time.Duration) (err error) {
	deadline := time.Now().Add(timeout)

	dailer := net.Dialer{Deadline: deadline}
	conn, err := dailer.Dial("unix", service)
	if err != nil {
		return
	}

	conn.SetDeadline(deadline)

	_, err = conn.Write([]byte{byte(request)})
	if err != nil {
		return
	}

	err = conn.Close()
	if err != nil {
		return
	}

	return
}
