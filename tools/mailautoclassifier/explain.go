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
	"os"
	"time"
)

const ps = string(os.PathSeparator)

var (
	explain             = true
	dictFilePath        = ".." + ps + ".." + ps + "data" + ps + "dict.txt"
	traningDataFilePath = "bayesian.data"
)

var testData = []struct {
	folder string
	class  string
}{
	{"good", string(analyzer.Good)},
	{"bad", string(analyzer.Bad)},
	{"neutral", "Neutral"},
}

func main() {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictFilePath)
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
			content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
			if err != nil {
				//ignore mail like Java Developer Day
			}

			class := anlz.Test(content)

			color := ""
			showInfo := false

			if item.class != class {
				totalError += 1
				showInfo = true
				color = "red+b"
				if class == analyzer.Neutral {
					totalNeutral += 1
				}
			} else if item.class == analyzer.Neutral {
				totalNeutral += 1
				showInfo = true
				color = "cyan+b"
			}

			if showInfo {
				msg := fmt.Sprintf("%s, %s\n", mail.Subject(), mailFilePath)
				fmt.Printf(ansi.Color(msg, color))
				if explain {
					if ba, ok := anlz.(*analyzer.BayesianAnalyzer); ok {
						fmt.Println(ba.Explain(content))
					}
				}
			}

		}
		elapsed := time.Now().Sub(startTime)
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
		fmt.Printf("TotalNum: %d, TotalError: %d, ErrRate: %f, TotalNeutral:%d, Confident:%f\n",
			totalNum, totalError, float64(totalError)/float64(totalNum), totalNeutral, float64(totalNum-totalNeutral)/float64(totalNum))
	}
}
