package main

import (
	"filter"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"os/signal"
	"postfix"
	"runtime"
	"time"
)

var (
	queuesFolder = "/var/spool/postfix/"
	//queuesFolder     = "data" + ps + "fakeQueues" + ps
	holdFolder       = queuesFolder + "hold"
	quarantineFolder = queuesFolder + "quarantine"
	incomingFolder   = queuesFolder + "incoming"

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictDataFilePath    = "data" + ps + "dict.data"

	allPass         = false
	subjectPrefixes = []string{"JWorld@TW新話題通知"}
	localDomain     = "javaworld.com.tw"
	cacheSize       = 100

	defaultMailFileFactory = &postfix.PostfixMailFileFactory{}
	numOfProcessor         = 100
	folderScanInterval     = time.Duration(1) * time.Second

	qmgrService  = "/var/spool/postfix/public/qmgr"
	flushTimeout = time.Duration(5) * time.Second

	logsFolder             = "log"
	metricLog              = logsFolder + ps + "metrics.log"
	writeMetricLogInterval = time.Duration(10) * time.Second
	spamdefenderLog        = logsFolder + ps + "spamdefender.log"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	startTime := time.Now()
	initDefaultLogger()

	log.Println("Starting daemon")
	startService()
	log.Printf("Daemon startup in %s", time.Since(startTime))

	waitForExit()
	log.Println("Stopping daemon")
	Shutdown()
	log.Println("Daemon stopped")
}

func initDefaultLogger() {
	sdl, err := os.OpenFile(spamdefenderLog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(sdl)

	AddShutdownHook(func() {
		sdl.Close()
	})
}

func startService() {

	defaultDestinationFilter := filter.NewDefaultDestinationFilter()
	contentInspectionFilter := filter.NewContentInspectionFilter(defaultDestinationFilter, allPass, traningDataFilePath, dictDataFilePath)
	cachingProxy := filter.NewCachingProxy(contentInspectionFilter, cacheSize)
	subjectPrefixMatchFilter := filter.NewSubjectPrefixMatchFilter(cachingProxy, subjectPrefixes)
	relayOnlyFilter := filter.NewRelayOnlyFilter(subjectPrefixMatchFilter, localDomain)

	paths := make(map[filter.Result]string)
	paths[filter.Incoming] = incomingFolder
	paths[filter.Quarantine] = quarantineFolder
	deliverFilter := filter.NewDeliverFilter(relayOnlyFilter, paths)

	handlerAdapter := filter.NewFileHandlerAdapter(deliverFilter, defaultMailFileFactory)

	flusher := NewPostfixFlusher(qmgrService, flushTimeout)
	dispatcher := NewPooledDispatcher(handlerAdapter, flusher, numOfProcessor)
	monitor := NewFolderMonitor(holdFolder, folderScanInterval, dispatcher)

	startMetric()
	monitor.Start()
	AddShutdownHook(func() {
		monitor.Stop()
	})

}

func startMetric() {
	ml, err := os.OpenFile(metricLog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	AddShutdownHook(func() {
		ml.Close()
	})

	go metrics.Log(metrics.DefaultRegistry, writeMetricLogInterval, log.New(ml, "", log.LstdFlags))
}

func waitForExit() {
	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)
	<-userInterrupt
}
