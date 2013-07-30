package general

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"errors"
	"fmt"
	"io"
	"net/url"
	"post"
	"strconv"
	"time"
)

func Parse(reader io.Reader) (newPost *post.Post, err error) {

	newPost = &post.Post{}
	currentIdx := 0
	parsers := []post.PartParser{&ReceiverParser{}, &SenderParser{}, &SubjectParser{}, &PostDateParser{}, &ContentParser{}, &LinkParser{}}
	bodyBlock := false

	z := html.NewTokenizer(reader)

loop:
	for {
		tokenType := z.Next()
		switch tokenType {
		case html.StartTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = true
			}
		case html.EndTagToken:
			if z.Token().DataAtom == atom.Body {
				bodyBlock = false
			}
		case html.TextToken:
			if bodyBlock {
				flow := parsers[currentIdx].Parse(newPost, z.Text())
				switch flow {
				case post.Next:
					if currentIdx < len(parsers) {
						currentIdx += 1
					}
				case post.Error:
					err = parsers[currentIdx].Err()
					break loop
				case post.Stop:
					break loop
				}
			}
		case html.ErrorToken:
			if z.Err() != io.EOF {
				err = z.Err()
			}
			break loop
		}
	}

	if currentIdx != len(parsers)-1 {
		err = errors.New("malformed email content")
	}

	return
}

type ReceiverParser struct {
	count int
	err   error
}

func (rp *ReceiverParser) Parse(newPost *post.Post, raw []byte) post.Flow {

	if rp.count > 2 {
		rp.err = errors.New("parse receiver failed: out of range")
		return post.Error
	}

	if !bytes.HasPrefix(raw, []byte("\nHi")) {
		rp.count += 1
		return post.Continue
	}

	fields := bytes.Split(raw, []byte(","))
	if len(fields) != 2 {
		rp.err = errors.New(fmt.Sprintf("parse receiver failed: %s", raw))
		return post.Error
	}

	newPost.Receiver = string(bytes.TrimSpace(fields[1]))
	return post.Next
}

func (rp *ReceiverParser) Err() error {
	return rp.err
}

type SenderParser struct {
	kind  int
	count int
}

func (sp *SenderParser) Parse(newPost *post.Post, raw []byte) post.Flow {

	// kind 0
	//JWorld@TW, 回覆通知 => sender at line 2

	// kind 1
	//JWorld@TW新話題通知:想請問這個Java Code那裡出錯了 => sender at line 3
	//JWorld@TW話題更新通知 => sender at line 3

	sp.count += 1

	if bytes.HasPrefix(raw, []byte("\n您訂閱的版面")) || bytes.HasPrefix(raw, []byte("\n您訂閱的主題")) {
		sp.kind = 1
	}

	if (sp.kind == 0 && sp.count < 2) || (sp.kind == 1 && sp.count < 3) {
		return post.Continue
	}

	newPost.Sender = string(bytes.TrimSpace(raw))
	return post.Next
}

func (sp *SenderParser) Err() error {
	return nil
}

type SubjectParser struct {
	ready bool
}

func (sp *SubjectParser) Parse(newPost *post.Post, raw []byte) post.Flow {

	if sp.ready {
		newPost.Subject = string(bytes.TrimSpace(raw))
		return post.Next
	}

	if bytes.HasPrefix(raw, []byte("\nSubject:")) {
		sp.ready = true
	}

	return post.Continue
}

func (sp *SubjectParser) Err() error {
	return nil
}

type PostDateParser struct {
	err error
}

func (pdp *PostDateParser) Parse(newPost *post.Post, raw []byte) post.Flow {
	dateStr := bytes.TrimPrefix(raw, []byte("\nDate:       "))
	//mail type: JWorld@TW新話題通知
	newPost.Date, pdp.err = time.Parse("2006-01-02 15:04", string(dateStr))
	if pdp.err != nil {
		//fallback:
		//mail type: JWorld@TW話題更新通知
		newPost.Date, pdp.err = time.Parse(time.UnixDate, string(dateStr))

		if pdp.err != nil {
			return post.Error
		}
	}

	return post.Next
}

func (pdp PostDateParser) Err() error {
	return pdp.err
}

type ContentParser struct {
	buf bytes.Buffer
}

func (cp *ContentParser) Parse(newPost *post.Post, raw []byte) post.Flow {
	if bytes.HasPrefix(raw, []byte("\n---------- ")) {
		newPost.Content = cp.buf.String()
		return post.Next
	}

	cp.buf.Write(raw)
	return post.Continue
}

func (cp *ContentParser) Err() error {
	return nil
}

type LinkParser struct {
	err error
}

func (lp *LinkParser) Parse(newPost *post.Post, raw []byte) post.Flow {

	if bytes.HasPrefix(raw, []byte("http://www.javaworld.com.tw/jute/post/view")) {
		newPost.Link = string(raw)

		var postLink *url.URL
		postLink, lp.err = url.Parse(newPost.Link)
		if lp.err != nil {
			return post.Error
		}

		values := postLink.Query()
		newPost.Id, lp.err = strconv.Atoi(values.Get("id"))
		if lp.err != nil {
			return post.Error
		}

		newPost.Bid, lp.err = strconv.Atoi(values.Get("bid"))
		if lp.err != nil {
			return post.Error
		}

		return post.Stop
	}

	return post.Continue
}

func (sp *LinkParser) Err() error {
	return sp.err
}
