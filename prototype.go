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

	numOfProcessor     = 100
	folderScanInterval = time.Duration(1) * time.Second

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictFilePath        = "data" + ps + "dict.txt"

	quit = make(chan struct{})
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sendOutOnly := mh.NewSendOutOnly(localDomain, incomingFolder)
	matchedSubject := mh.NewMatchedSubject(subjectPrefix, incomingFolder)
	contentInspection := mh.NewContentInspection(allPass, quarantineFolder, traningDataFilePath, dictFilePath)
	defaultDestination := mh.NewDefaultDestination(incomingFolder)

	handlerChain := mh.NewHandlerChain(sendOutOnly, matchedSubject, contentInspection, defaultDestination)
	handlerAdapter := mh.NewFileHandlerAdapter(handlerChain, &mailfile.POP3MailFileFactory{})
	dispatcher := service.NewDispatcher(handlerAdapter, numOfProcessor)

	monitor := service.NewFolderMonitor(holdFolder, folderScanInterval, dispatcher)
	monitor.Start()

	<-quit
	monitor.Stop()
}
