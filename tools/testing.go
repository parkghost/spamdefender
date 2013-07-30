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
	confident           = 0.01
	dictDataFilePath    = filepath.Join("..", "data", "dict.data")
	traningDataFilePath = filepath.Join("..", "data", "bayesian.data")
)

var testData = []struct {
	folder string
	class  string
}{
	{filepath.Join("..", "data", "test", "good"), analyzer.Good},
	{filepath.Join("..", "data", "test", "bad"), analyzer.Bad},
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
			fileutil.Humanize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))
	}
}
