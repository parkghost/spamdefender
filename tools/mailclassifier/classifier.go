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

	totalNum, goods, bads, neutrals, badformats := 0, 0, 0, 0, 0

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
		content, err := htmlutil.ExtractText(htmlText, htmlutil.BannerRemover("----------", 0, 1))
		mail.Close()

		color := ""
		moveTo := ""

		if err != nil {
			badformats += 1
			moveTo = "badformat" + ps + fi.Name()
		} else {
			class := anlz.Test(content)

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
		}

		if !dryRun {
			common.CopyFile(filePath, moveTo)
		}

		msg := fmt.Sprintf("%s %s\n", mail.Subject(), moveTo)
		fmt.Printf(ansi.Color(msg, color))
	}
	fmt.Printf("TotalNum: %d, goods: %d, bads: %d, neutrals: %d, badformat:%d\n", totalNum, goods, bads, neutrals, badformats)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
