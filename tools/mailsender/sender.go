package main

import (
	"common"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	folder      = filepath.Join("..", "mailfetecher", "mailbox")
	recipient   = "alice@labs.brandonc.me"
	from        = "brandon@labs.brandonc.me"
	user        = "*"
	passoword   = "*"
	smtpServer  = "mail.labs.brandonc.me"
	port        = 25
	concurrency = 50 //postfix: smtpd_client_connection_count_limit
)

func main() {
	log.Printf("Mail all mails from %s to %s", folder, recipient)
	totalNum := 0
	var totalSize int64

	startTime := time.Now()
	fis, err := ioutil.ReadDir(folder)
	checkErr(err)

	auth := smtp.PlainAuth(
		"",
		user,
		passoword,
		smtpServer,
	)

	clients, err := NewSmtpClient(smtpServer+":"+strconv.Itoa(port), auth, concurrency)
	checkErr(err)
	defer clients.Close()

	wg := &sync.WaitGroup{}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		totalSize += fi.Size()
		totalNum += 1

		filePath := filepath.Join(folder, fi.Name())
		rawBody, err := ioutil.ReadFile(filePath)
		checkErr(err)

		wg.Add(1)
		go func() {
			err = clients.SendMail(
				from,
				[]string{recipient},
				rawBody,
			)
			checkErr(err)
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Now().Sub(startTime)
	fmt.Printf("TotalNum: %d, Elapsed: %s, TPS(Mail): %f, TPS(FileSize): %s\n",
		totalNum,
		time.Now().Sub(startTime),
		float64(totalNum)/(float64(elapsed)/float64(time.Second)),
		common.HumanReadableSize(uint64(float64(totalSize)/(float64(elapsed)/float64(time.Second)))))

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type SmtpClient struct {
	pool   chan *smtp.Client
	mu     *sync.Mutex
	closed bool
}

func (s *SmtpClient) SendMail(from string, to []string, msg []byte) (err error) {
	c := <-s.pool
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return
		}
	}
	w, err := c.Data()
	if err != nil {
		return
	}
	_, err = w.Write(msg)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	s.pool <- c
	return
}

func (s *SmtpClient) Close() {
	close(s.pool)
	for c := range s.pool {
		err := c.Quit()
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func NewSmtpClient(addr string, a smtp.Auth, concurrency int) (client *SmtpClient, err error) {
	pool := make(chan *smtp.Client, concurrency)
	for i := 0; i < concurrency; i++ {
		var c *smtp.Client
		c, err = smtp.Dial(addr)
		if err != nil {
			return
		}
		if err = c.Hello("localhost"); err != nil {
			return
		}
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err = c.StartTLS(nil); err != nil {
				return
			}
		}

		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return
			}
		}

		pool <- c
	}

	client = &SmtpClient{pool, &sync.Mutex{}, false}

	return
}
