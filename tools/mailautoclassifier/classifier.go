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
)

const ps = string(os.PathSeparator)

var (
	cutset = ":;=<>"

	confident           = 0.01
	dryRun              = false
	mailbox             = ".." + ps + "mailfetecher" + ps + "mailbox"
	dictFilePath        = ".." + ps + ".." + ps + "data" + ps + "dict.txt"
	traningDataFilePath = "bayesian.data"
)

func main() {
	anlz, err := analyzer.NewAnalyzer(traningDataFilePath, dictFilePath)
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

		class, score := anlz.Test(content)

		color := ""
		if math.Abs(score[analyzer.Good]/score[analyzer.Bad]-1) < confident {
			color = "cyan+b"
			neutrals += 1
			if !dryRun {
				common.CopyFile(filePath, "neutral"+ps+fi.Name())
			}
		} else {
			switch class {
			case analyzer.Good:
				color = "green+b"
				goods += 1
				if !dryRun {
					common.CopyFile(filePath, "good"+ps+fi.Name())
				}
			case analyzer.Bad:
				color = "red+b"
				bads += 1
				if !dryRun {
					common.CopyFile(filePath, "bad"+ps+fi.Name())
				}
			}
		}

		totalNum += 1

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
