package mailfile

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type Testdata struct {
	source      string
	plainSource string
}

type plainSourceReader struct {
	reader *bufio.Reader
}

func (p *plainSourceReader) ReadLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(line, "\t\r\n"), nil
}

func (p *plainSourceReader) ReadRestString() (string, error) {

	var buf bytes.Buffer
	_, err := p.reader.WriteTo(&buf)
	if err != nil {
		return "", err
	}

	return strings.Trim(buf.String(), "\t\r\n"), nil
}

func NewPlainSourceReader(reader io.Reader) *plainSourceReader {
	return &plainSourceReader{bufio.NewReader(reader)}
}

var pop3Testdata = Testdata{
	"testdata" + string(os.PathSeparator) + "pop3",
	"testdata" + string(os.PathSeparator) + "pop3_plain",
}

func TestPOP3RetrieveSubject(t *testing.T) {
	mail := NewPOP3Mail(pop3Testdata.source)
	f, err := os.Open(pop3Testdata.plainSource)
	defer f.Close()

	reader := NewPlainSourceReader(f)
	subject, err := reader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}

	if subject != mail.Subject() {
		t.Fatalf("expected mail subject is %s, got %s", subject, mail.Subject())
	}
}

func TestPOP3RetrieveContent(t *testing.T) {
	mail := NewPOP3Mail(pop3Testdata.source)
	f, err := os.Open(pop3Testdata.plainSource)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	reader := NewPlainSourceReader(f)
	_, err = reader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}

	exptectedBodyContent, err := reader.ReadRestString()
	if err != nil {
		t.Fatal(err)
	}

	gotBodyConent := strings.Trim(mail.Content(), "\t\r\n")

	if exptectedBodyContent != gotBodyConent {
		t.Fatalf("expected mail content is \n%v\n, got %v", exptectedBodyContent, gotBodyConent)
	}
}
