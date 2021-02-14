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
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/lidar", handleLidar)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	log.Printf("request from %s: %s %q", req.RemoteAddr, req.Method, req.URL)
	fmt.Fprintln(w, "hello from roverd!")
}

func handleLidar(w http.ResponseWriter, req *http.Request) {
	lidarCmd := req.URL.Query().Get("lidar")

	action := "did nothing with lidar-scan"

	if lidarCmd == "0" {
		err := StopLidar()
		if err != nil {
			log.Print(err)
		} else {
			action = "stopped lidar-scan"
		}

	} else if lidarCmd == "1" {
		pid, err := StartLidar()
		if err != nil {
			log.Print(err)
		} else {
			action = fmt.Sprintf("started lidar-scan (pid %d)", pid)
		}
	}

	log.Print(action)
	fmt.Fprintln(w, action)
}
