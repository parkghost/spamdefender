package main

import (
	"bufio"
	"goseg"
	"io/ioutil"
	"log"
	"mailfile"
	"mailpost"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var allData = []struct {
	folder string
	output string
}{
	{"good", "goodtermfreq.txt"},
	{"bad", "badtermfreq.txt"},
}

var (
	cutset = ":;=<>"

	termMinLength    = 1
	topTermRatio     = 0.90
	dictDataFilePath = filepath.Join("..", "..", "data", "dict.data")
)

func main() {
	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range allData {
		fis, err := ioutil.ReadDir(item.folder)

		termFreq := make(map[string]int)

		if err != nil {
			log.Fatal(err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			filePath := filepath.Join(item.folder, fi.Name())
			mail := mailfile.NewPOP3Mail(filePath)
			if err = mail.Parse(); err != nil {
				log.Fatal(err)
			}

			post, err := mailpost.Parse(mail)
			mail.Close()
			if err != nil {
				log.Fatalf("Err: %v, Mail:%s", err, mail.Path())
			}

			words := tokenizer.Cut([]rune(post.Subject + " " + post.Content))

			for _, word := range words {
				key := strings.Trim(word, cutset)
				if len(key) > 1 {
					termFreq[key] = termFreq[key] + 1
				}
			}
		}

		pairList := sortMapByValue(termFreq)

		output, err := os.OpenFile(item.output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		writer := bufio.NewWriter(output)

		offset := int(float64(len(pairList)) * topTermRatio)

		for _, pair := range pairList[offset:] {
			if len([]rune(pair.Key)) > termMinLength {
				writer.WriteString(pair.Key + " " + strconv.Itoa(pair.Value) + "\n")
			}
		}
		writer.Flush()
		output.Close()

	}
}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i += 1
	}
	sort.Sort(p)
	return p
}
