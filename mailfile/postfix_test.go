package mailfile

import (
	"os"
	"strings"
	"testing"
)

var postfixTestdata = Testdata{
	"testdata" + string(os.PathSeparator) + "postfix",
	"testdata" + string(os.PathSeparator) + "postfix_plain",
}

func TestPostfixRetrieveSubject(t *testing.T) {
	mail := NewPostfixMail(postfixTestdata.source)
	f, err := os.Open(postfixTestdata.plainSource)
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

func TestPostfixRetrieveContent(t *testing.T) {
	mail := NewPostfixMail(postfixTestdata.source)
	f, err := os.Open(postfixTestdata.plainSource)
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
