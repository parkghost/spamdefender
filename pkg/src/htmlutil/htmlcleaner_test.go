package htmlutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func BannerRemover(lineSeparator string, skipTop int, skipBottom int) func(string) (string, error) {
	return func(text string) (string, error) {

		lines := strings.Split(text, "\n")

		var pos []int
		for no, line := range lines {
			if strings.TrimRight(line, " ") == lineSeparator {
				pos = append(pos, no)
			}
		}

		if len(pos) == 0 {
			return text, errors.New(fmt.Sprintf("html: cannot detect lineSeparator: %s", lineSeparator))
		}

		if len(pos) != 3 {
			return text, errors.New("html: malformed mail content")
		}

		top, bottom := pos[skipTop], pos[len(pos)-skipBottom-1]
		return strings.Join(lines[top+1:bottom-1], "\n"), nil
	}
}

var testdata = struct {
	original string
	expected string
}{
	filepath.Join("testdata", "original"),
	filepath.Join("testdata", "expected"),
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
