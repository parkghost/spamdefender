package filter

import (
	"mailfile"
)

const (
	None = iota
	Incoming
	Quarantine
)

type Result int

type Filter interface {
	Filter(mailfile.Mail) Result
}
