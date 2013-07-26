package main

import (
	"analyzer"
	"common"
	"fmt"
	"github.com/mgutz/ansi"
	"htmlutil"
	"io/ioutil"
	"log"
	"mailfile"
	"os"
	"time"
)

const ps = string(os.PathSeparator)

var (
	confident           = 0.01
	dictDataFilePath    = ".." + ps + "data" + ps + "dict.data"
	traningDataFilePath = ".." + ps + "data" + ps + "bayesian.data"
)

var testData = []struct {
	folder string
	class  string
}{
	{".." + ps + "data" + ps + "test" + ps + "good", string(analyzer.Good)},
	{".." + ps + "data" + ps + "test" + ps + "bad", string(analyzer.Bad)},
}

func main() {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range testData {
		log.Printf("Testing %s", item.folder)

		totalNum, totalError, totalNeutral := 0, 0, 0
		var totalSize int64

		fis, err := ioutil.ReadDir(item.folder)

		if err != nil {
			log.Fatal(err)
		}

		startTime := time.Now()
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			totalSize += fi.Size()
			totalNum += 1

			mailFilePath := item.folder + ps + fi.Name()
			mail := mailfile.NewPOP3Mail(mailFilePath)
			if err = mail.Parse(); err != nil {
				log.Fatal(err)
			}

			htmlText := mail.Content()
			content, _ := htmlutil.ExtractText(htmlText, htmlutil.BannerRemover("----------", 0, 1))
			mail.Close()

			class := anlz.Test(content)

			switch {
			case analyzer.Neutral == class:
				totalNeutral += 1
				fmt.Println(ansi.Color(mailFilePath, "cyan+b"))

			case item.class != class:
				totalError += 1
				fmt.Println(ansi.Color(mailFilePath, "red+b"))
			}

		}

		elapsed := time.Now().Sub(startTime)

		fmt.Printf("TotalNum: %d, TotalError: %d, ErrRate: %f, TotalNeutral:%d, Confident:%f\n",
			totalNum, totalError, float64(totalError)/float64(totalNum), totalNeutral, float64(totalNum-totalNeutral)/float64(totalNum))
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
	}
}
