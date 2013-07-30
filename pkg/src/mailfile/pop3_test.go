package mailfile

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type Testdata struct {
	source      string
	plainSource string
}

var pop3Testdata = Testdata{
	filepath.Join("testdata", "pop3"),
	filepath.Join("testdata", "pop3_plaintext"),
}

func TestPOP3RetrieveSubject(t *testing.T) {
	pop3mail := NewPOP3Mail(pop3Testdata.source)
	defer pop3mail.Close()
	if err := pop3mail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := NewPlainTextMail(pop3Testdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if plainTextMail.Subject() != pop3mail.Subject() {
		t.Fatalf("expected %s, got %s", plainTextMail.Subject(), pop3mail.Subject())
	}
}

func TestPOP3RetrieveContent(t *testing.T) {
	pop3mail := NewPOP3Mail(pop3Testdata.source)
	defer pop3mail.Close()
	if err := pop3mail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := NewPlainTextMail(pop3Testdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	actual, err := readContentString(pop3mail.Content())
	if err != nil {
		t.Fatal(err)
	}
	actual = strings.TrimSpace(actual)

	expected, err := readContentString(plainTextMail.Content())
	if err != nil {
		t.Fatal(err)
	}
	expected = strings.TrimSpace(expected)

	if expected != actual {
		t.Fatalf("expected \n%v\n, got %v", expected, actual)
	}
}

func readContentString(reader io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestPOP3RetrieveFrom(t *testing.T) {
	pop3mail := NewPOP3Mail(pop3Testdata.source)
	defer pop3mail.Close()
	if err := pop3mail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := NewPlainTextMail(pop3Testdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(plainTextMail.From(), pop3mail.From()) {
		t.Fatalf("expected \n%v\n, got %v", plainTextMail.From(), pop3mail.From())
	}
}

func TestPOP3RetrieveTo(t *testing.T) {
	pop3mail := NewPOP3Mail(pop3Testdata.source)
	defer pop3mail.Close()
	if err := pop3mail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := NewPlainTextMail(pop3Testdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(plainTextMail.To(), pop3mail.To()) {
		t.Fatalf("expected \n%v\n, got %v", plainTextMail.From(), pop3mail.From())
	}
}
