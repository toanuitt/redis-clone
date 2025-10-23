package main

import (
	"os"
	"os/signal"
	"redis-clone/internal/server"
	"sync"
	"syscall"
)

func main() {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)
	go server.RunIoMultiplexingServer(&wg)
	go server.WaitForSignal(&wg, signals)
	wg.Wait()

}
