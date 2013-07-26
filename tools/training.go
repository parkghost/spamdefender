package main

import (
	"common"
	"fmt"
	"github.com/parkghost/bayesian"
	"goseg"
	"htmlutil"
	"io/ioutil"
	"log"
	"mailfile"
	"os"
	"time"
)

const ps = string(os.PathSeparator)

var (
	Good   bayesian.Class = "Good"
	Bad    bayesian.Class = "Bad"
	cutset                = "1234567890:;=<>"
)

var (
	dictDataFilePath = ".." + ps + "data" + ps + "dict.data"
	output           = ".." + ps + "data" + ps + "bayesian.data"
)

var trainingData = []struct {
	folder string
	class  bayesian.Class
}{
	{".." + ps + "data" + ps + "training" + ps + "good", Good},
	{".." + ps + "data" + ps + "training" + ps + "bad", Bad},
}

func main() {
	classifier := bayesian.NewClassifier(Good, Bad)
	tokenizer, err := goseg.NewTokenizerFromFile(dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range trainingData {
		log.Printf("Traning %s data", item.class)
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

			filePath := item.folder + string(os.PathSeparator) + fi.Name()
			mail := mailfile.NewPOP3Mail(filePath)
			if err = mail.Parse(); err != nil {
				log.Fatal(err)
			}

			htmlText := mail.Content()
			content, _ := htmlutil.ExtractText(htmlText, htmlutil.BannerRemover("----------", 0, 1))
			mail.Close()

			words := common.Normalize(tokenizer.Cut([]rune(content)), cutset)
			classifier.Learn(words, item.class)
			totalNum += 1
		}

		elapsed := time.Now().Sub(startTime)
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
	}

	classifier.WriteToFile(output)
}
