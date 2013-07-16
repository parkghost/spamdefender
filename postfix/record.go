package postfix

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	space = 0x20
)

// postfix-2.10.1/src/global/rec_type.h
// postfix-2.10.1/src/global/record.c
const (
	VmailerExtraOffs = iota
	Postfix10DataOffset
	Postfix10RecipientCount
	Postfix21qmgrFlags
	Postfix24ContentLength
)

type RecordReader struct {
	reader *bufio.Reader
	buf    *bytes.Buffer
	parsed bool
	err    error
}

func (r *RecordReader) Read(p []byte) (n int, err error) {
	if !r.parsed {
		err = r.parse()
		if err != nil {
			if err != io.EOF {
				r.err = err
				return
			}
		}
	}

	if r.err != nil {
		return 0, r.err
	}

	return r.buf.Read(p)
}

func (r *RecordReader) parse() (err error) {
	if r.parsed {
		return r.err
	}
	r.parsed = true
	n, pos := 0, 0

	// parse first record group (envelope information)
	contentSizes, n, err := parseFirstRecord(r.reader)
	if err != nil {
		return
	}
	pos += n

	// ignore rest records
	n, err = skip(r.reader, contentSizes[Postfix10DataOffset]-pos)
	if err != nil {
		return
	}
	pos += n

	// parse second record group (mail content)
	for {
		var data []byte
		data, n, err = readRawRecord(r.reader)
		if err != nil {
			return
		}
		pos += n

		_, err = r.buf.Write(data)
		if err != nil {
			return
		}

		_, err = r.buf.WriteString("\n")
		if err != nil {
			return
		}

		if pos > contentSizes[Postfix10DataOffset]+contentSizes[Postfix24ContentLength] {
			break
		}
	}

	// ignore third record group

	return
}

func parseFirstRecord(reader *bufio.Reader) (contentSizes []int, length int, err error) {
	contentSizes = make([]int, 5)

	// skip REC_TYPE_SIZE
	recTypes, err := reader.ReadBytes(space)
	if err != nil {
		return
	}
	length += len(recTypes)

	// The record at the start of the queue file specifies the message content size
	// Vmailer extra offs - data offs
	// Postfix 1.0 data offset
	// Postfix 1.0 recipient count
	// Postfix 2.1 qmgr flags
	// Postfix 2.4 content length

	// read envelope message content sizes
	for i := 0; i < len(contentSizes); i++ {
		// REC_TYPE_SIZE_FORMAT	"%15ld %15ld %15ld %15ld %15ld"
		var span []byte
		n := 0
		if i == 0 {
			span, n, err = read(reader, 15-1) // reader.ReadBytes already read one character(space)
		} else {
			span, n, err = read(reader, 16)
		}
		length += n
		if err != nil {
			return
		}

		contentSizes[i], err = strconv.Atoi(strings.TrimSpace(string(span)))
		if err != nil {
			return
		}
	}

	return
}

func read(reader *bufio.Reader, size int) (buf []byte, length int, err error) {
	buf = make([]byte, size)
	for {
		var n int
		n, err = reader.Read(buf[length:])
		length += n
		if length == size {
			return
		}

		if err != nil {
			break
		}
	}
	return
}

func skip(reader *bufio.Reader, size int) (int, error) {
	var i int
	for i < size {
		_, err := reader.ReadByte()
		if err != nil {
			return i, err
		}
		i++
	}

	return i, nil
}

func readRawRecord(reader *bufio.Reader) (data []byte, length int, err error) {
	// parse rec_type and data length
	meta, err := reader.Peek(5)
	if err != nil {
		return
	}

	dataLength, offset, err := parseBase128Int(meta, 1)
	if err != nil {
		return
	}

	// read data
	n, err := skip(reader, offset)
	if err != nil {
		return
	}
	length += n

	data, n, err = read(reader, dataLength)
	if err != nil {
		return
	}
	length += n

	return
}

func parseBase128Int(bytes []byte, initOffset int) (ret, offset int, err error) {
	offset = initOffset
	for shifted := 0; offset < len(bytes); shifted += 7 {
		if shifted > 32 {
			err = errors.New("base 128 integer too large")
			return
		}
		b := bytes[offset]
		ret |= int(b&0x7f) << uint(shifted)
		offset++
		if b&0x80 == 0 {
			return
		}
	}
	err = errors.New("truncated base 128 integer")
	return
}

func NewRecordReader(reader io.Reader) *RecordReader {
	return &RecordReader{reader: bufio.NewReader(reader), buf: &bytes.Buffer{}}
}
