package goseg

import (
	"os"
	"testing"
)

const ps = string(os.PathSeparator)

func TestTrieAddChild(t *testing.T) {
	root := &TrieNode{}

	text := "ebdca"
	for _, c := range []rune(text) {
		root.addChild(c)
	}

	expected := "[a:[], b:[], c:[], d:[], e:[]]"
	actual := root.String()
	if expected != actual {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestTrieLookup(t *testing.T) {
	root := &TrieNode{}

	text := "ebdca"
	for _, c := range []rune(text) {
		root.addChild(c)
	}
	for _, c := range []rune(text) {
		if node := root.Lookup(c); node == nil {
			t.Fatalf("expected %v, got %v", string(c), nil)
		} else {
			if c != node.Char {
				t.Fatalf("expected %v, got %v", string(c), string(node.Char))
			}
		}
	}

	if node := root.Lookup('g'); node != nil {
		t.Fatalf("expected nil, but got %v", node)
	}
}

func TestTrieAddString(t *testing.T) {
	root := &TrieNode{}

	root.AddString("ab")
	root.AddString("bc")
	root.AddString("cd")

	expected := "[a:[b:[]], b:[c:[]], c:[d:[]]]"
	actual := root.String()
	if expected != actual {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
