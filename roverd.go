package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/sevlyar/go-daemon"
)

var (
	sigChan = make(chan os.Signal, 1)
)

// To terminate the daemon use:
// kill `cat roverd.pid`
func main() {
	ctx := &daemon.Context{
		PidFileName: "roverd.pid",
		PidFilePerm: 0644,
		LogFileName: "roverd.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	daemon, err := ctx.Reborn()
	if err != nil {
		log.Fatal("unable to run: ", err)
	}
	if daemon != nil {
		return
	}
	defer ctx.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	signal.Notify(sigChan, os.Interrupt)
	go listenInterrupt()

	serveHTTP()
}

func listenInterrupt() {
	for sig := range sigChan {
		log.Printf("dameon stopped by %s", sig.String())
		os.Exit(0)
	}
}

func serveHTTP() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("request from %s: %s %q", req.RemoteAddr, req.Method, req.URL)
	fmt.Fprintln(w, "hello from roverd!")
}
