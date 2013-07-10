package main

import (
	"fmt"
	"github.com/parkghost/bayesian"
	"github.com/parkghost/spamdefender/analyzer/goseg"
	"github.com/parkghost/spamdefender/common"
	"io/ioutil"
	"log"
	"os"
)

const ps = string(os.PathSeparator)

var (
	Good   bayesian.Class = "Good"
	Bad    bayesian.Class = "Bad"
	cutset                = "1234567890:;=<>"
)

var (
	dictDataFilePath = ".." + ps + ".." + ps + "data" + ps + "dict.data"
	output           = "bayesian.data"
)

func main() {
	classifier := bayesian.NewClassifier(Good, Bad)

	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
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