package goseg

/*
forked from https://github.com/fxsjy/goseg
removed artificial neural network for reducing memeory usage and speed up
add Tokenizer struct for configurable and serializable
TrieTree use struct instead of map for reducing memeory usage(slightly lower lookup performance)
*/

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type tuple struct {
	freq float32
	pos  int
}

type TrieNode struct {
	Children []*TrieNode
	Char     rune
	Last     bool
}

func (tn *TrieNode) String() string {
	buf := bytes.NewBufferString("[")
	for i, node := range tn.Children {
		if i == 0 {
			buf.WriteRune(node.Char)
			buf.WriteString(":" + node.String())
		} else {
			buf.WriteString(", ")
			buf.WriteRune(node.Char)
			buf.WriteString(":" + node.String())
		}

	}
	buf.WriteString("]")
	return buf.String()
}

func (tn *TrieNode) Lookup(Char rune) *TrieNode {
	n := len(tn.Children)
	if n == 0 {
		return nil
	}
	// binary search
	l, r := 0, n-1
	for {
		m := (l + r) / 2
		c := tn.Children[m].Char
		if c == Char {
			return tn.Children[m]
		} else if c > Char {
			r = m - 1
		} else {
			l = m + 1
		}

		if l > r {
			return nil
		}
	}
}

func (tn *TrieNode) AddString(word string) *TrieNode {
	ptr := tn
	for _, c := range []rune(word) {
		if node := ptr.Lookup(c); node != nil {
			ptr = node
		} else {
			ptr = ptr.addChild(c)
		}
	}
	ptr.Last = true
	return ptr
}

func (tn *TrieNode) addChild(Char rune) *TrieNode {
	node := &TrieNode{Children: make([]*TrieNode, 0), Char: Char}

	n := len(tn.Children)
	if n == 0 {
		tn.Children = append(tn.Children, node)
	} else {
		// add node to sorted slice
		for i := 0; i < n; i++ {
			child := tn.Children[i]
			if Char < child.Char {
				newChildren := make([]*TrieNode, 0, len(tn.Children)+1)
				newChildren = append(newChildren, tn.Children[:i]...)
				newChildren = append(newChildren, node)
				newChildren = append(newChildren, tn.Children[i:]...)
				tn.Children = newChildren
				break
			}

			if i == n-1 {
				tn.Children = append(tn.Children, node)
			}
		}
	}
	return node
}

type Tokenizer struct {
	trie       *TrieNode
	freq_table map[string]float32
	min_freq   float32
}

func (tk *Tokenizer) Cut(sentence []rune) []string {
	i := 0
	j := 0
	N := len(sentence)
	words := make([]string, 0)
	for i < N {
		c := sentence[i]
		j = i
		if isEng(c) {
			for i < N && isEng(sentence[i]) {
				i++
			}
			words = append(words, string(sentence[j:i]))
		} else if isHan(c) {
			for i < N && isHan(sentence[i]) {
				i++
			}
			tmp := tk.cut_DAG(sentence[j:i])
			words = append(words, tmp...)

		} else {
			for i < N && !isEng(sentence[i]) && !isHan(sentence[i]) {
				i++
			}
		}
	}
	return words
}

func isHan(c rune) bool {
	if c >= 19968 && c <= 40869 {
		return true
	}
	return false
}

func isEng(c rune) bool {
	if c >= 48 && c <= 122 {
		return true
	}
	return false
}

func (tk *Tokenizer) cut_DAG(sentence []rune) []string {
	N := len(sentence)
	i, j := 0, 0
	p := tk.trie
	DAG := make(map[int][]int)
	for i < N {
		c := sentence[j]
		if p.Lookup(c) != nil {
			p = p.Lookup(c)
			if p.Last {
				if DAG[i] == nil {
					DAG[i] = make([]int, 0)
				}
				DAG[i] = append(DAG[i], j)
			}
			j++
			if j >= N {
				i++
				j = i
				p = tk.trie
			}
		} else {
			p = tk.trie
			i++
			j = i
		}
	}

	for i := 0; i < N; i++ {
		if DAG[i] == nil {
			DAG[i] = []int{i}
		}
	}

	route := make(map[int]*tuple)
	tk.calc(sentence, DAG, 0, route)
	x := 0
	words := make([]string, 0)
	for x < N {
		y := route[x].pos + 1
		l_word := sentence[x:y]
		x = y
		words = append(words, string(l_word))
	}
	return words
}

func (tk *Tokenizer) calc(sentence []rune, DAG map[int][]int, idx int, route map[int]*tuple) {
	N := len(sentence)
	route[N] = &tuple{freq: 1.0, pos: -1}

	for idx := N - 1; idx >= 0; idx-- {
		best := &tuple{freq: -1, pos: -1}
		next := DAG[idx]
		for _, x := range next {
			candidate := route[x+1]
			word_freq := tk.freq_table[string(sentence[idx:x+1])]
			if word_freq == 0 { //smooth
				word_freq = tk.min_freq
			}
			prod := candidate.freq * word_freq
			if prod > best.freq {
				best.freq = prod
				best.pos = x
			}
		}
		route[idx] = best
	}

}

func NewTokenizer(dictionaries ...string) (*Tokenizer, error) {
	trie := &TrieNode{}
	freq_table := make(map[string]float32)
	var min_freq float32 = 0.0

	for _, dict := range dictionaries {
		file, err := os.Open(dict)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var content []byte
		content, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		var content_str string
		content_str = string(content)
		var lines []string
		lines = strings.Split(content_str, "\n")
		for _, line := range lines {
			var word string = ""
			var freq int = 0
			tup := strings.Split(line, " ")
			word = tup[0]
			freq, _ = strconv.Atoi(tup[1])
			trie.AddString(word)
			freq_table[word] = float32(freq)
		}
	}

	min_freq, freq_table = normalize(freq_table)

	return &Tokenizer{trie, freq_table, min_freq}, nil
}

func normalize(d map[string]float32) (float32, map[string]float32) {
	new_d := make(map[string]float32)
	var sum float32 = 0.0
	for _, v := range d {
		sum += v
	}
	var _min float32 = 1.0
	for k, v := range d {
		t := v / sum
		new_d[k] = t
		if t < _min {
			_min = t
		}
	}
	return _min, new_d
}

type serializableTokenizer struct {
	Trie      *TrieNode
	FreqTable map[string]float32
	MinFreq   float32
}

func NewTokenizerFromFile(name string) (c *Tokenizer, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewTokenizerFromReader(file)
}

func NewTokenizerFromReader(r io.Reader) (c *Tokenizer, err error) {
	dec := gob.NewDecoder(r)
	w := new(serializableTokenizer)
	err = dec.Decode(w)
	if err != nil {
		return nil, err
	}

	return &Tokenizer{w.Trie, w.FreqTable, w.MinFreq}, nil
}

func (tk *Tokenizer) WriteToFile(name string) (err error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return tk.WriteTo(file)
}

func (tk *Tokenizer) WriteTo(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	err = enc.Encode(&serializableTokenizer{tk.trie, tk.freq_table, tk.min_freq})
	return
}
