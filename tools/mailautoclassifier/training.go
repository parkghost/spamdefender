package main

import (
	"fmt"
	"github.com/jbrukh/bayesian"
	"io/ioutil"
	"log"
	"os"
	"spamdefender/analyzer/goseg"
	"spamdefender/common"
)

const ps = string(os.PathSeparator)

var (
	Good   bayesian.Class = "Good"
	Bad    bayesian.Class = "Bad"
	cutset                = "1234567890:;=<>"
)

var (
	dictFilePath = ".." + ps + ".." + ps + "data" + ps + "dict.txt"
	output       = "bayesian.data"
)

func main() {
	classifier := bayesian.NewClassifier(Good, Bad)

	tokenizer, err := goseg.NewTokenizer(dictFilePath)
	checkErr(err)

	goodWords, err := getWords("goodwords.txt", tokenizer)
	checkErr(err)
	fmt.Printf("Normalized Good Words:\n%s\n", goodWords)

	fmt.Println("")

	badWords, err := getWords("badwords.txt", tokenizer)
	checkErr(err)
	fmt.Printf("Normalized Bad Words:\n%s\n", badWords)

	classifier.Learn(goodWords, Good)
	classifier.Learn(badWords, Bad)

	classifier.WriteToFile(output)
}

func getWords(filePath string, tokenizer *goseg.Tokenizer) ([]string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	words := tokenizer.Cut([]rune(string(bytes)))
	return common.Normalize(words, cutset), nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
