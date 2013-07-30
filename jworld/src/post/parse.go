package post

import (
	"io"
	"time"
)

const (
	Next = iota
	Continue
	Stop
	Error
)

type Post struct {
	Id       int
	Bid      int
	Receiver string
	Sender   string
	Subject  string
	Date     time.Time
	Content  string
	Link     string
}

type PostParser interface {
	Parse(reader io.Reader) *Post
}

type Flow int

type PartParser interface {
	Parse(post *Post, raw []byte) Flow
	Err() error
}
