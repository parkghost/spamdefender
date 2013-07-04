package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	pop3 "github.com/d3xter/GoPOP3"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"spamdefender/common"
	"strconv"
	"strings"
)

var (
	username      = "your email account"
	password      = "your password "
	serverAddress = "pop3.live.com:995"
	secure        = true
	destination   = "mailbox"
	mailFrom      = "webmaster@mail.javaworld.com.tw"
	start         = 0
)

// TODO: retry download when failure
// TODO: multi-client for speed up

func main() {

	log.Printf("Connect to %s\n", serverAddress)
	var client *pop3.Client
	var dialErr error
	if secure {
		client, dialErr = DialTLS(serverAddress)
	} else {
		client, dialErr = pop3.Dial(serverAddress)
	}
	checkErr(dialErr)

	log.Println("Authenticating")
	plainAuth := pop3.CreatePlainAuthentication(username, password)
	authErr := client.Authenticate(plainAuth)
	checkErr(authErr)

	log.Println("Fetching mailbox stat")
	cmdErr := PrintMailBoxStat(client)
	checkErr(cmdErr)

	log.Println("Downloading mails")
	err := ScanMailbox(client)
	checkErr(err)

	log.Println("Finished to download all of mails")

}

func DialTLS(addr string) (client *pop3.Client, err error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host := addr[:strings.Index(addr, ":")]

	return pop3.NewClient(conn, host)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PrintMailBoxStat(client *pop3.Client) error {
	mailCount, mailBoxSize, err := client.GetStatus()
	if err != nil {
		return err
	}

	fmt.Printf("Account:%s, MailCount:%d, MailBoxSize:%s\n", username, mailCount, common.HumanReadableSize(uint64(mailBoxSize)))
	return nil
}

func ScanMailbox(client *pop3.Client) error {
	rawMailList, respErr := client.GetRawMailList()
	checkErr(respErr)

	scanner := bufio.NewScanner(strings.NewReader(rawMailList))
	scanner.Split(ScanFirstWordAtLine)

	lastNo := 0
	for scanner.Scan() {
		if scanner.Err() == nil {
			index, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatalf("Stop scan due to Err:%v, Token:%s", err, scanner.Text())
			}

			if index >= start {
				err = CheckAndDownloadMailContent(client, index)
				if err != nil {
					log.Fatalf("Download mail failure, Err:%v, No:%d, ", err, index)
				}
				lastNo = index
			}
		} else {
			log.Fatalf("Stop scan due to Err:%v, lastNo:%d", scanner.Err(), lastNo)
			break
		}
	}
	return nil
}

func ScanFirstWordAtLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		_, t, e := bufio.ScanWords(data, atEOF)
		return i + 1, t, e
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		_, t, e := bufio.ScanWords(data, atEOF)
		return len(data), t, e
	}
	// Request more data.
	return 0, nil, nil
}

func CheckAndDownloadMailContent(client *pop3.Client, index int) (err error) {
	filePath := destination + string(os.PathSeparator) + strconv.Itoa(index)
	if _, errFileStat := os.Stat(filePath); errFileStat != nil { // mail file not found

		log.Printf("Downloading mail, No:%d", index)
		// TODO: store binary file or handle non-utf8 encoding
		// THINK: check content integrity
		var mailContent string
		mailContent, err = client.GetRawMail(index)
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
			log.Printf("store mail content to file:%s", filePath)
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
