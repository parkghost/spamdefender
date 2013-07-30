package main

import (
	"analyzer"
	"fileutil"
	"fmt"
	"github.com/mgutz/ansi"
	"io/ioutil"
	"log"
	"mailfile"
	"mailpost"
	"path/filepath"
	"time"
)

var (
	explain             = true
	dictDataFilePath    = filepath.Join("..", "..", "data", "dict.data")
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

			mailFilePath := filepath.Join(item.folder, fi.Name())
			mail := mailfile.NewPOP3Mail(mailFilePath)
			if err = mail.Parse(); err != nil {
				log.Fatal(err)
			}

			post, err := mailpost.Parse(mail)
			mail.Close()
			if err != nil {
				log.Fatalf("Err: %v, Mail:%s", err, mail.Path())
			}

			class := anlz.Test(post.Subject + " " + post.Content)

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
					fmt.Println(anlz.Explain(post.Subject + " " + post.Content))
				}
			}

		}
		elapsed := time.Now().Sub(startTime)

		fmt.Printf("TotalNum: %d, TotalError: %d, ErrRate: %f, TotalNeutral:%d, Confident:%f\n",
			totalNum, totalError, float64(totalError)/float64(totalNum), totalNeutral, float64(totalNum-totalNeutral)/float64(totalNum))
		fmt.Printf("Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
			time.Now().Sub(startTime),
			float64(totalNum)/(float64(elapsed)/float64(time.Second)),
			fileutil.Humanize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
	}
}
