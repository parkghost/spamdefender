package goseg

import (
	"path/filepath"
	"testing"
)

func TestTrieNodeLookup(t *testing.T) {
	text := "ebdca"
	root := &TrieNode{}

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

}

func TestTrieNodeFind(t *testing.T) {
	data := []string{"ab", "bc", "cd"}
	root := &TrieNode{}

	for _, w := range data {
		root.AddString(w)
		if !root.Find(w) {
			t.Fatalf("%s not found", w)
		}
	}
}

func BenchmarkCut(b *testing.B) {
	b.StopTimer()
	text := `卡納維洛在報告中指出，他已經成功用PEG接合過老鼠和狗的脊髓，但目前若要應用到人類的頭顱移植手術，他還必須面對2大困難：資金和道德。卡納維洛認為，若要2年內達到人頭移植的目標，他至少需要3千萬美金的資金投入研究，以拯救那些為肌肉萎縮症、癱瘓、器官衰竭，甚至癌症的重症病患`
	tk, err := NewTokenizer(filepath.Join("..", "..", "..", "conf", "dict.txt"))
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tk.Cut([]rune(text))
	}
}

func BenchmarkNewTokenizerFromFile(b *testing.B) {
	b.ReportAllocs()
	dictDataFile := filepath.Join("..", "..", "..", "data", "dict.data")
	for i := 0; i < b.N; i++ {
		_, _ = NewTokenizerFromFile(dictDataFile)
	}
}

func BenchmarkNewTokenizer(b *testing.B) {
	b.ReportAllocs()
	dictFile := filepath.Join("..", "..", "..", "conf", "dict.txt")
	for i := 0; i < b.N; i++ {
		_, _ = NewTokenizer(dictFile)
	}
}
