package main

import (
	"bufio"
	"github.com/parkghost/spamdefender/analyzer/goseg"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const ps = string(os.PathSeparator)

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
	dictDataFilePath = ".." + ps + ".." + ps + "data" + ps + "dict.data"
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

			filePath := item.folder + string(os.PathSeparator) + fi.Name()
			mail := mailfile.NewPOP3Mail(filePath)
			if err = mail.Parse(); err != nil {
				log.Fatal(err)
			}

			htmlText := mail.Content()
			content, _ := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
			mail.Close()

			words := tokenizer.Cut([]rune(content))

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
