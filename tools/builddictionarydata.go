package main

import (
	"goseg"
	"log"
	"path/filepath"
)

var (
	dictFilePath     = filepath.Join("..", "conf", "dict.txt")
	dictDataFilePath = filepath.Join("..", "data", "dict.data")
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
