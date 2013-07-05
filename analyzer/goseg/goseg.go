package goseg

/*
forked from https://github.com/fxsjy/goseg
removed artificial neural network for speedup
*/
import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type trie_node struct {
	lookup map[rune]*trie_node
}

type tuple struct {
	freq float32
	pos  int
}

func new_trie_node() *trie_node {
	new_node := &trie_node{}
	new_node.lookup = make(map[rune]*trie_node)
	return new_node
}

func add_string(trie *trie_node, s string) {
	var ptr *trie_node
	ptr = trie
	for _, c := range s {
		if ptr.lookup[c] == nil {
			ptr.lookup[c] = new_trie_node()
		}
		ptr = ptr.lookup[c]
	}
	ptr.lookup[0] = new_trie_node()
}

type Tokenizer struct {
	trie       *trie_node
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
		if p.lookup[c] != nil {
			p = p.lookup[c]
			if p.lookup[0] != nil {
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
			//fmt.Println(idx,x+1,N)
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
	trie := new_trie_node()
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
			add_string(trie, word)
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
