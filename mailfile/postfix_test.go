package mailfile

import (
	"io/ioutil"
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
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err = mail.Parse(); err != nil {
		t.Fatal(err)
	}
	defer mail.Close()

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
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err = mail.Parse(); err != nil {
		t.Fatal(err)
	}
	defer mail.Close()

	reader := NewPlainSourceReader(f)
	_, err = reader.ReadLine()
	if err != nil {
		t.Fatal(err)
	}

	exptectedBodyContent, err := reader.ReadRestString()
	if err != nil {
		t.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(mail.Content())
	if err != nil {
		t.Fatal(err)
	}
	gotBodyConent := strings.Trim(string(bodyBytes), "\t\r\n")

	if exptectedBodyContent != gotBodyConent {
		t.Fatalf("expected mail content is \n%v\n, got %v", exptectedBodyContent, gotBodyConent)
	}
}
