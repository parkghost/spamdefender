package html

import (
	"io/ioutil"
	"os"
	"testing"
)

const ps = string(os.PathSeparator)

var testdata = struct {
	original string
	expected string
}{
	"testdata" + ps + "original",
	"testdata" + ps + "expected",
}

func TestExtractText(t *testing.T) {
	raw, err := ioutil.ReadFile(testdata.original)
	if err != nil {
		t.Fatal(err)
	}
	originalText := string(raw)

	raw, err = ioutil.ReadFile(testdata.expected)
	if err != nil {
		t.Fatal(err)
	}
	expectedText := string(raw)

	newText, err := ExtractText(originalText, BannerRemover("----------", 0, 1))
	if err != nil {
		t.Fatal(err)
	}

	if expectedText != newText {
		t.Fatalf("expected %s, got %s", expectedText, newText)
	}
}
