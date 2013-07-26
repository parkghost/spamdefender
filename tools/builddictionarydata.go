package main

import (
	"goseg"
	"log"
	"os"
)

const ps = string(os.PathSeparator)

var (
	dictFilePath     = ".." + ps + "conf" + ps + "dict.txt"
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
