package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"./services"

	"github.com/ordishs/gocore"
)

var (
	activeCoins, _ = gocore.Config().Get("activeCoins")
	coins          = strings.Split(activeCoins, ",")
)

func main() {
	stats := gocore.Config().Stats()
	log.Printf("STATS\n%s\nVERSION\n-------\n%s (%s)\n\n", stats, version, commit)

	// setup signal catching
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		s := <-signalChan

		log.Printf("Received signal: %s", s)
		appCleanup()
		os.Exit(1)
	}()

	start()
}

func appCleanup() {
	log.Println("shareprocessor shutting dowm...")
}

func start() {
	services.ConnectToZMQ()

	waitCh := make(chan bool)
	<-waitCh
}
