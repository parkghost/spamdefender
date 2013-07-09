package main

import (
	"github.com/parkghost/spamdefender/analyzer/goseg"
	"log"
	"os"
)

const ps = string(os.PathSeparator)

var (
	dictFilePath     = ".." + ps + "data" + ps + "dict.txt"
	dictDataFilePath = ".." + ps + "data" + ps + "dict.data"
)

func main() {
	tokenizer, err := goseg.NewTokenizer(dictFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = tokenizer.WriteToFile(dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
