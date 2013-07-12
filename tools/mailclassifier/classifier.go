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
)

const ps = string(os.PathSeparator)

var (
	cutset = ":;=<>"

	confident           = 0.01
	dryRun              = false
	mailbox             = ".." + ps + "mailfetecher" + ps + "mailbox"
	dictDataFilePath    = ".." + ps + ".." + ps + "data" + ps + "dict.data"
	traningDataFilePath = "bayesian.data"
)

func main() {
	anlz, err := analyzer.NewBayesianAnalyzer(traningDataFilePath, dictDataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	totalNum, goods, bads, neutrals := 0, 0, 0, 0

	fis, err := ioutil.ReadDir(mailbox)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		totalNum += 1

		filePath := mailbox + ps + fi.Name()
		mail := mailfile.NewPOP3Mail(filePath)
		if err = mail.Parse(); err != nil {
			log.Fatal(err)
		}

		htmlText := mail.Content()
		content, err := html.ExtractText(htmlText, html.BannerRemover("----------", 0, 1))
		if err != nil {
			fmt.Println(err)
		}
		mail.Close()

		class := anlz.Test(content)

		color := ""
		moveTo := ""

		switch class {
		case analyzer.Neutral:
			neutrals += 1
			color = "cyan+b"
			moveTo = "neutral" + ps + fi.Name()
		case analyzer.Good:
			goods += 1
			color = "green+b"
			moveTo = "good" + ps + fi.Name()
		case analyzer.Bad:
			bads += 1
			color = "red+b"
			moveTo = "bad" + ps + fi.Name()
		}

		if !dryRun {
			common.CopyFile(filePath, moveTo)
		}

		msg := fmt.Sprintf("%s %s\n", mail.Subject(), filePath)
		fmt.Printf(ansi.Color(msg, color))
	}
	fmt.Printf("TotalNum: %d, goods: %d, bads: %d, neutrals: %d\n", totalNum, goods, bads, neutrals)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
