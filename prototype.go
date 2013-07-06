package main

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	mh "github.com/parkghost/spamdefender/service/mail"
	"os"
	"runtime"
	"time"
)

const ps = string(os.PathSeparator)

var (
	allPass       = false
	localDomain   = "javaworld.com.tw"
	subjectPrefix = "JWorld@TW新話題通知"

	//baseFolder = "/var/spool/postfix/"
	baseFolder       = "fakeQueues" + ps
	holdFolder       = baseFolder + "hold"
	quarantineFolder = baseFolder + "quarantine"
	incomingFolder   = baseFolder + "incoming"

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictFilePath        = "data" + ps + "dict.txt"

	quit = make(chan struct{})
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	contentInspection := mh.NewContentInspection(allPass, quarantineFolder, traningDataFilePath, dictFilePath)
	sendOutOnly := mh.NewSendOutOnly(localDomain, incomingFolder)
	matchedSubject := mh.NewMatchedSubject(subjectPrefix, incomingFolder)
	defaultDestination := mh.NewDefaultDestination(incomingFolder)

	handlerChain := mh.NewHandlerChain(sendOutOnly, matchedSubject, contentInspection, defaultDestination)
	handler := mh.NewFileHandlerAdapter(handlerChain, &mailfile.POP3MailFileFactory{})
	dispatcher := service.NewDispatcher(handler, 100)

	monitor := service.NewFolderMonitor(holdFolder, time.Duration(1)*time.Second, dispatcher)
	monitor.Start()

	<-quit
	monitor.Stop()
}
