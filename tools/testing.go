package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/parkghost/spamdefender/analyzer"
	"github.com/parkghost/spamdefender/common"
	"github.com/parkghost/spamdefender/html"
	"github.com/parkghost/spamdefender/mailfile"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

const ps = string(os.PathSeparator)

var (
	confident           = 0.01
	dictFilePath        = ".." + ps + "data" + ps + "dict.txt"
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

	anlz, err := analyzer.NewAnalyzer(traningDataFilePath, dictFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range testData {
		log.Printf("Testing %s", item.folder)

		totalNum, totalError, totalConfident := 0, 0, 0
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
			content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
			if err != nil {
				//ignore mail like Java Developer Day
			}

			class, score := anlz.Test(content)

			if math.Abs(score[analyzer.Good]/score[analyzer.Bad]-1) > confident {
				totalConfident += 1
			} else {
				msg := fmt.Sprintf("%s, %f\n", mailFilePath, score[analyzer.Good]/score[analyzer.Bad])
				fmt.Printf(ansi.Color(msg, "cyan+b"))
			}

			if item.class != string(class) {
				totalError += 1
				msg := fmt.Sprintf("%s, %f\n", mailFilePath, score[analyzer.Good]/score[analyzer.Bad])
				fmt.Printf(ansi.Color(msg, "red+b"))
			}
		}

		elapsed := time.Now().Sub(startTime)
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
		fmt.Printf("TotalNum: %d, TotalError: %d, ErrRate: %f, TotalConfident:%d, ConfidentRate:%f\n",
			totalNum, totalError, float64(totalError)/float64(totalNum), totalConfident, float64(totalConfident)/float64(totalNum))
	}
}
