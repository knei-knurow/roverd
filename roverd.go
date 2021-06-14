package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

var signals = make(chan os.Signal, 1)

func init() {
	log.SetFlags(0)
	log.SetPrefix("roverd: ")
}

func main() {
	ctx := &daemon.Context{
		PidFileName: "roverd.pid",
		PidFilePerm: 0644,
		LogFileName: "roverd.log",
		LogFilePerm: 0644,
		WorkDir:     "./",
		Umask:       027,
	}
	defer ctx.Release()

	_, err := ctx.Reborn()
	if err != nil {
		log.Fatalln("failed to reborn context:", err)
	}

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go listenSignals()
}

func listenSignals() {
	for sig := range signals {
		log.Printf("dameon stopped by %s\n", sig.String())
		os.Exit(0)
	}
}
