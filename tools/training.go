package main

import (
	"common"
	"fmt"
	"github.com/parkghost/bayesian"
	"goseg"
	"io/ioutil"
	"log"
	"mailfile"
	"mailpost"
	"path/filepath"
	"time"
)

var (
	Good   bayesian.Class = "Good"
	Bad    bayesian.Class = "Bad"
	cutset                = "1234567890:;=<>"
)

var (
	dictDataFilePath = filepath.Join("..", "data", "dict.data")
	output           = filepath.Join("..", "data", "bayesian.data")
)

var trainingData = []struct {
	folder string
	class  bayesian.Class
}{
	{filepath.Join("..", "data", "training", "good"), Good},
	{filepath.Join("..", "data", "training", "bad"), Bad},
}

func main() {
	classifier := bayesian.NewClassifier(Good, Bad)
	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range trainingData {
		log.Printf("Traning %s", item.folder)
		totalNum := 0
		var totalSize int64

		startTime := time.Now()
		fis, err := ioutil.ReadDir(item.folder)

		if err != nil {
			log.Fatal(err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			totalSize += fi.Size()

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

			words := common.Normalize(tokenizer.Cut([]rune(post.Subject+" "+post.Content)), cutset)
			classifier.Learn(words, item.class)
			totalNum += 1
		}

		elapsed := time.Now().Sub(startTime)

		fmt.Printf("TotalNum: %d\n", totalNum)
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))

	}

	classifier.WriteToFile(output)
}
