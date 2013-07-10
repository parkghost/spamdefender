package main

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	"github.com/parkghost/spamdefender/service/filter"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const ps = string(os.PathSeparator)

var (
	allPass         = false
	localDomain     = "javaworld.com.tw"
	subjectPrefixes = []string{"JWorld@TW新話題通知", "JWorld@TW話題更新通知", "JWorld@TW新文章通知"}

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
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	startTime := time.Now()
	log.Println("Starting daemon")

	defaultDestinationFilter := filter.NewDefaultDestinationFilter(incomingFolder)
	contentInspectionFilter := filter.NewContentInspectionFilter(defaultDestinationFilter, allPass, quarantineFolder, traningDataFilePath, dictDataFilePath)
	subjectPrefixMatchFilter := filter.NewSubjectPrefixMatchFilter(contentInspectionFilter, subjectPrefixes, incomingFolder)
	relayOnlyFilter := filter.NewRelayOnlyFilter(subjectPrefixMatchFilter, localDomain, incomingFolder)
	cachingFilter := filter.NewCachingFilter(relayOnlyFilter, cacheSize)

	handlerAdapter := filter.NewFileHandlerAdapter(cachingFilter, &mailfile.POP3MailFileFactory{})
	dispatcher := service.NewPooledDispatcher(handlerAdapter, numOfProcessor)

	monitor := service.NewFolderMonitor(holdFolder, folderScanInterval, dispatcher)
	log.Printf("Daemon startup in %s", time.Since(startTime))
	monitor.Start()

	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)
	<-userInterrupt

	log.Println("Stopping daemon")
	monitor.Stop()
	log.Println("Daemon stopped")
}
