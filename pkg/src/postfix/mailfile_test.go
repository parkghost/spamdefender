package postfix

import (
	"io"
	"io/ioutil"
	"mailfile"
	"os"
	"reflect"
	"strings"
	"testing"
)

type Testdata struct {
	source      string
	plainSource string
}

var postfixTestdata = Testdata{
	"testdata" + string(os.PathSeparator) + "postfix",
	"testdata" + string(os.PathSeparator) + "postfix_plaintext",
}

func TestPostfixRetrieveSubject(t *testing.T) {
	postfixMail := NewPostfixMail(postfixTestdata.source)
	defer postfixMail.Close()
	if err := postfixMail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := mailfile.NewPlainTextMail(postfixTestdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if plainTextMail.Subject() != postfixMail.Subject() {
		t.Fatalf("expected mail subject is %s, got %s", plainTextMail.Subject(), postfixMail.Subject())
	}
}

func TestPostfixRetrieveContent(t *testing.T) {
	postfixMail := NewPostfixMail(postfixTestdata.source)
	defer postfixMail.Close()
	if err := postfixMail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := mailfile.NewPlainTextMail(postfixTestdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	actual, err := readContentString(postfixMail.Content())
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
		t.Fatalf("expected mail content is \n%v\n, got %v", expected, actual)
	}
}

func readContentString(reader io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestPostfixRetrieveFrom(t *testing.T) {
	postfixMail := NewPostfixMail(postfixTestdata.source)
	defer postfixMail.Close()
	if err := postfixMail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := mailfile.NewPlainTextMail(postfixTestdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(plainTextMail.From(), postfixMail.From()) {
		t.Fatalf("expected \n%v\n, got %v", plainTextMail.From(), postfixMail.From())
	}
}

func TestPostfixRetrieveTo(t *testing.T) {
	postfixMail := NewPostfixMail(postfixTestdata.source)
	defer postfixMail.Close()
	if err := postfixMail.Parse(); err != nil {
		t.Fatal(err)
	}

	plainTextMail := mailfile.NewPlainTextMail(postfixTestdata.plainSource)
	defer plainTextMail.Close()
	if err := plainTextMail.Parse(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(plainTextMail.To(), postfixMail.To()) {
		t.Fatalf("expected \n%v\n, got %v", plainTextMail.From(), postfixMail.From())
	}
}
