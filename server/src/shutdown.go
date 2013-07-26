package main

import (
	"sync"
)

var (
	shutdownFnList []func()
	shutdownMux    sync.Mutex
)

func AddShutdownHook(fun func()) {
	shutdownMux.Lock()
	shutdownFnList = append(shutdownFnList, fun)
	shutdownMux.Unlock()
}

func Shutdown() {
	shutdownMux.Lock()
	for _, fn := range shutdownFnList {
		fn()
	}
	shutdownMux.Unlock()
}
