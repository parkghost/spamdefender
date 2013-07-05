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
	explain             = true
	dictFilePath        = ".." + ps + ".." + ps + "data" + ps + "dict.txt"
	traningDataFilePath = "bayesian.data"
)

var testData = []struct {
	folder string
	test   string
}{
	{"good", "good"},
	{"bad", "bad"},
	{"neutral", "neutral"},
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

			score, pass := anlz.Test(content)
			testConfident := math.Abs(score[0]/score[1] - 1)

			if testConfident < confident {
				msg := fmt.Sprintf("%s, %s, %f\n", mail.Subject(), mailFilePath, testConfident)
				fmt.Printf(ansi.Color(msg, "cyan+b"))
				if explain {
					fmt.Println(anlz.Explain(content))
				}
			} else {
				totalConfident += 1
				if (item.test == "good" && !pass) ||
					(item.test == "bad" && pass) ||
					(item.test == "neutral") {
					totalError += 1
					msg := fmt.Sprintf("%s, %s, %f\n", mail.Subject(), mailFilePath, testConfident)
					fmt.Printf(ansi.Color(msg, "red+b"))
					if explain {
						fmt.Println(anlz.Explain(content))
					}
				}
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
