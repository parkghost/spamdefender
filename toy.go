package main

import (
	"os"
	"spamdefender/app"
	mh "spamdefender/app/mailhandler"
	"time"
)

const ps = string(os.PathSeparator)

var (
	allPass       = false
	localDomain   = "labs.brandonc.me"
	subjectPrefix = "JWorld@TW新話題通知"

	baseFolder = "/var/spool/postfix" + ps
	//baseFolder       = "testdata" + ps + "fakeQueues" + ps
	holdFolder       = baseFolder + "hold"
	quarantineFolder = baseFolder + "quarantine"
	incomingFolder   = baseFolder + "incoming"

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictFilePath        = "data" + ps + "dict.txt"

	quit = make(chan struct{})
)

func main() {

	contentInspection := mh.NewContentInspection(allPass, quarantineFolder, traningDataFilePath, dictFilePath)
	sendOutOnly := mh.NewSendOutOnly(localDomain, incomingFolder)
	matchedSubject := mh.NewMatchedSubject(subjectPrefix, incomingFolder)
	finalDestination := mh.NewFinalDestination(incomingFolder)

	handlerChain := mh.NewHandlerChain(sendOutOnly, matchedSubject, contentInspection, finalDestination)
	handler := mh.NewMailHandlerAdapter(handlerChain)

	monitor := app.NewFolderMonitor(holdFolder, time.Duration(1)*time.Second, handler)
	monitor.Start()

	<-quit
	monitor.Stop()
}
