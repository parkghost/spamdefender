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
)

var (
	cutset = ":;=<>"

	confident           = 0.01
	dryRun              = false
	mailbox             = filepath.Join("..", "mailfetecher", "mailbox")
	dictDataFilePath    = filepath.Join("..", "..", "data", "dict.data")
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

		filePath := filepath.Join(mailbox, fi.Name())
		mail := mailfile.NewPOP3Mail(filePath)
		if err = mail.Parse(); err != nil {
			log.Fatal(err)
		}

		post, err := mailpost.Parse(mail)
		mail.Close()
		if err != nil {
			log.Printf("Err: %v, Mail:%s\n", err, mail.Path())
		}

		color := ""
		moveTo := ""

		if err != nil {
			badformats += 1
			moveTo = filepath.Join("badformat", fi.Name())
		} else {
			class := anlz.Test(post.Subject + " " + post.Content)

			switch class {
			case analyzer.Neutral:
				neutrals += 1
				color = "cyan+b"
				moveTo = filepath.Join("neutral", fi.Name())
			case analyzer.Good:
				goods += 1
				color = "green+b"
				moveTo = filepath.Join("good", fi.Name())
			case analyzer.Bad:
				bads += 1
				color = "red+b"
				moveTo = filepath.Join("bad", fi.Name())
			}
		}

		if !dryRun {
			fileutil.CopyFile(filePath, moveTo)
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
