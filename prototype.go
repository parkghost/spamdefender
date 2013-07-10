package main

import (
	"github.com/parkghost/spamdefender/mailfile"
	"github.com/parkghost/spamdefender/service"
	"github.com/parkghost/spamdefender/service/filter"
	metrics "github.com/rcrowley/go-metrics"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const ps = string(os.PathSeparator)

var (
	//queuesFolder = "/var/spool/postfix/"
	queuesFolder     = "fakeQueues" + ps
	holdFolder       = queuesFolder + "hold"
	quarantineFolder = queuesFolder + "quarantine"
	incomingFolder   = queuesFolder + "incoming"

	traningDataFilePath = "data" + ps + "bayesian.data"
	dictDataFilePath    = "data" + ps + "dict.data"

	allPass         = false
	subjectPrefixes = []string{"JWorld@TW新話題通知", "JWorld@TW話題更新通知", "JWorld@TW新文章通知"}
	localDomain     = "javaworld.com.tw"
	cacheSize       = 100

	numOfProcessor     = 100
	folderScanInterval = time.Duration(1) * time.Second

	logsFolder             = "logs"
	metricLog              = logsFolder + ps + "metric.log"
	writeMetricLogInterval = time.Duration(10) * time.Second
	spamdefenderLog        = logsFolder + ps + "spamdefender.log"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sdl, err := os.OpenFile(spamdefenderLog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.Close()
	log.SetOutput(sdl)

	startTime := time.Now()
	log.Println("Starting daemon")

	defaultDestinationFilter := filter.NewDefaultDestinationFilter(incomingFolder)
	contentInspectionFilter := filter.NewContentInspectionFilter(defaultDestinationFilter, allPass, quarantineFolder, traningDataFilePath, dictDataFilePath)
	subjectPrefixMatchFilter := filter.NewSubjectPrefixMatchFilter(contentInspectionFilter, subjectPrefixes, incomingFolder)
	relayOnlyFilter := filter.NewRelayOnlyFilter(subjectPrefixMatchFilter, localDomain, incomingFolder)
	cachingFilter := filter.NewCachingFilter(relayOnlyFilter, cacheSize)
	deliverFilter := filter.NewDeliverFilter(cachingFilter)

	handlerAdapter := filter.NewFileHandlerAdapter(deliverFilter, &mailfile.POP3MailFileFactory{})
	dispatcher := service.NewPooledDispatcher(handlerAdapter, numOfProcessor)
	monitor := service.NewFolderMonitor(holdFolder, folderScanInterval, dispatcher)
	log.Printf("Daemon startup in %s", time.Since(startTime))
	monitor.Start()

	ml, err := os.OpenFile(metricLog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer ml.Close()
	go metrics.Log(metrics.DefaultRegistry, writeMetricLogInterval, log.New(ml, "metrics: ", log.Lmicroseconds))

	userInterrupt := make(chan os.Signal, 1)
	signal.Notify(userInterrupt, os.Interrupt)
	<-userInterrupt

	log.Println("Stopping daemon")
	monitor.Stop()
	log.Println("Daemon stopped")
}
