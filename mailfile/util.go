package mailfile

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

func DecodeRFC2047String(s string) (string, error) {

	lines := strings.Split(s, " ")

	buf := bytes.NewBufferString("")
	for _, line := range lines {
		if strings.HasPrefix(line, "=?") {
			word, err := DecodeRFC2047Word(line)
			if err != nil {
				return "", err
			}
			buf.WriteString(word)
		} else {
			buf.WriteString(" " + line)
		}
	}

	return buf.String(), nil
}

func DecodeRFC2047Word(s string) (string, error) {
	fields := strings.Split(s, "?")
	if len(fields) != 5 || fields[0] != "=" || fields[4] != "=" {
		return "", errors.New("mail: address not RFC 2047 encoded")
	}
	charset, enc := strings.ToLower(fields[1]), strings.ToLower(fields[2])
	if charset != "iso-8859-1" && charset != "utf-8" {
		return "", fmt.Errorf("mail: charset not supported: %q", charset)
	}

	in := bytes.NewBufferString(fields[3])
	var r io.Reader
	switch enc {
	case "b":
		r = base64.NewDecoder(base64.StdEncoding, in)
	case "q":
		r = qDecoder{r: in}
	default:
		return "", fmt.Errorf("mail: RFC 2047 encoding not supported: %q", enc)
	}

	dec, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	switch charset {
	case "iso-8859-1":
		b := new(bytes.Buffer)
		for _, c := range dec {
			b.WriteRune(rune(c))
		}
		return b.String(), nil
	case "utf-8":
		return string(dec), nil
	}
	panic("unreachable")
}

type qDecoder struct {
	r       io.Reader
	scratch [2]byte
}

func (qd qDecoder) Read(p []byte) (n int, err error) {
	// This method writes at most one byte into p.
	if len(p) == 0 {
		return 0, nil
	}
	if _, err := qd.r.Read(qd.scratch[:1]); err != nil {
		return 0, err
	}
	switch c := qd.scratch[0]; {
	case c == '=':
		if _, err := io.ReadFull(qd.r, qd.scratch[:2]); err != nil {
			return 0, err
		}
		x, err := strconv.ParseInt(string(qd.scratch[:2]), 16, 64)
		if err != nil {
			return 0, fmt.Errorf("mail: invalid RFC 2047 encoding: %q", qd.scratch[:2])
		}
		p[0] = byte(x)
	case c == '_':
		p[0] = ' '
	default:
		p[0] = c
	}
	return 1, nil
}
