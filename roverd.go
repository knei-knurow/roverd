package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/sevlyar/go-daemon"
)

// To terminate the daemon use:
//  kill `cat roverd.pid`
func main() {
	ctx := &daemon.Context{
		PidFileName: "roverd.pid",
		PidFilePerm: 0644,
		LogFileName: "roverd.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
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

	serveHTTP()
}

func serveHTTP() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from %s: %s %q", r.RemoteAddr, r.Method, r.URL)
	fmt.Fprintf(w, "hello from roverd: %q", html.EscapeString(r.URL.Path))
}
