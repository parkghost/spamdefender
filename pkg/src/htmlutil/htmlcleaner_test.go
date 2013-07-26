package htmlutil

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
	originalFile, err := os.Open(testdata.original)
	if err != nil {
		t.Fatal(err)
	}

	expectedBytes, err := ioutil.ReadFile(testdata.expected)
	if err != nil {
		t.Fatal(err)
	}
	expectedText := string(expectedBytes)

	originalText, err := ExtractText(originalFile, BannerRemover("----------", 0, 1))
	if err != nil {
		t.Fatal(err)
	}

	if expectedText != originalText {
		t.Fatalf("expected %s, got %s", expectedText, originalText)
	}
}
