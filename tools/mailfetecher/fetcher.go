package main

import (
	"bytes"
	"fileutil"
	"fmt"
	pop3 "github.com/bytbox/go-pop3"
	"github.com/parkghost/pkg/net/mail"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	username      = "your email account"
	password      = "your password "
	serverAddress = "pop3.live.com:995"
	secure        = true
	destination   = "mailbox"
	mailFrom      = "webmaster@mail.javaworld.com.tw"
	offset        = 0
)

// TODO: retry download on failure

func main() {

	log.Printf("Connect to %s\n", serverAddress)

	var client *pop3.Client
	var dialErr error
	if secure {
		client, dialErr = pop3.DialTLS(serverAddress)
	} else {
		client, dialErr = pop3.Dial(serverAddress)
	}
	checkErr(dialErr)

	authErr := client.Auth(username, password)
	checkErr(authErr)

	log.Println("Fetching mailbox stat")
	cmdErr := PrintMailBoxStat(client)
	checkErr(cmdErr)

	log.Println("Downloading mails")
	err := ScanMailbox(client)
	checkErr(err)

	log.Println("Finished to download all of mails")

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PrintMailBoxStat(client *pop3.Client) error {
	mailCount, mailBoxSize, err := client.Stat()
	if err != nil {
		return err
	}

	fmt.Printf("Account:%s, MailCount:%d, MailBoxSize:%s\n", username, mailCount, fileutil.Humanize(uint64(mailBoxSize)))
	return nil
}

func ScanMailbox(client *pop3.Client) error {
	msgs, _, err := client.ListAll()
	checkErr(err)

	for i := offset; i < len(msgs); i++ {
		err = CheckAndDownloadMailContent(client, msgs[i])
		if err != nil {
			log.Fatalf("Download mail failure, Err:%v, Mail:%d, ", err, msgs[i])
		}
	}

	return nil
}

func CheckAndDownloadMailContent(client *pop3.Client, index int) (err error) {
	filePath := filepath.Join(destination, strconv.Itoa(index))
	if _, errFileStat := os.Stat(filePath); errFileStat != nil { // mail file not found

		log.Printf("Downloading Mail:%d", index)
		// TODO: handle non-utf8 encoding
		// THINK: check content integrity
		var mailContent string
		mailContent, err = client.Retr(index)
		if err != nil {
			return
		}

		//find newline at head of file
		var lastNewline = 0
		for i, b := range mailContent {
			if b != '\n' {
				lastNewline = i
				break
			}
		}
		rawMailContent := []byte(mailContent)[lastNewline:]

		// it's safe to ignore parse error
		fromAddress, _ := getFromAddress(rawMailContent)

		if fromAddress == mailFrom {
			log.Printf("Store mail content to file:%s", filePath)
			err = ioutil.WriteFile(filePath, rawMailContent, 0644)
			if err != nil {
				return
			}
		}

	}
	return
}

func getFromAddress(rawMailContent []byte) (string, error) {
	var message *mail.Message
	message, err := mail.ReadMessage(bytes.NewReader(rawMailContent))
	if err != nil {
		return "", err
	}

	address, err := mail.ParseAddress(message.Header.Get("From"))
	if err != nil {
		log.Printf("parse from address failure, Err:%v, Raw:%s", err, message.Header.Get("From"))
		return "", err
	}

	return address.Address, nil
}
