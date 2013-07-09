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
	cacheSize          = 100
	folderScanInterval = time.Duration(1) * time.Second

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictDataFilePath    = "data" + ps + "dict.data"

	quit = make(chan struct{})
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	defaultDestination := mh.NewDestination(incomingFolder)
	contentInspection := mh.NewContentInspection(defaultDestination, allPass, quarantineFolder, traningDataFilePath, dictDataFilePath)
	//subjectPrefixMatch := mh.NewSubjectPrefixMatch(contentInspection, subjectPrefix, incomingFolder)
	sendOutOnly := mh.NewSendOutOnly(contentInspection, localDomain, incomingFolder)
	cache := mh.NewCache(sendOutOnly, cacheSize)

	handlerAdapter := mh.NewFileHandlerAdapter(cache, &mailfile.POP3MailFileFactory{})
	dispatcher := service.NewPooledDispatcher(handlerAdapter, numOfProcessor)

	monitor := service.NewFolderMonitor(holdFolder, folderScanInterval, dispatcher)
	monitor.Start()

	<-quit
	monitor.Stop()
}
