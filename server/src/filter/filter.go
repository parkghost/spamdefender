package filter

import (
	"mailfile"
	"os"
)

const (
	ps = string(os.PathSeparator)

	None       = Result(0)
	Incoming   = Result(1)
	Quarantine = Result(2)
)

type Result int

type Filter interface {
	Filter(mailfile.Mail) Result
}
