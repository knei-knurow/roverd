package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/knei-knurow/roverd/handlers/move"
	"github.com/knei-knurow/roverd/modules/motors"
	"github.com/knei-knurow/roverd/modules/servos"
	"github.com/knei-knurow/roverd/ports"
	"github.com/tarm/serial"
)

var (
	// Host on which the roverd will listen for requests. Usually empty or "localhost".
	host string

	// Port on which the roverd will listen for requests.
	port string

	// Whether to log extensive output.
	verbose bool
)

var (
	// Port with the device controlling motors.
	movePort ports.Serial

	// Port with the device controlling rangefinder.
	// rangefinderPort sercom.Serial
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("roverd: ")

	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.Parse()

	if verbose {
		motors.Verbose = true
		servos.Verbose = true
	}

	env := Env{}
	env.Load()
	log.Printf("loaded env vars: %v\n", env)

	host = env.listenHost
	port = env.listenPort
	movePortName := env.movePort
	movePortBaud, err := strconv.Atoi(env.movePortBaud)
	if err != nil {
		log.Fatalf("cannot read baud rate: %v\n", err)
	}

	config := &serial.Config{
		Name: movePortName,
		Baud: movePortBaud,
	}

	p, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("cannot open port %s: %v\n", movePortName, err)
	}

	movePort = ports.SerialPort{ReadWriteCloser: p}

	motors.Port = movePort
	servos.Port = movePort
}

func main() {
	addr := host + ":" + port

	log.Println("listening on", addr)
	serveHTTP(addr)
}

func serveHTTP(addr string) {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/move", handleMove)
	http.ListenAndServe(addr, nil)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	log.Printf("request from %s: %s %q", req.RemoteAddr, req.Method, req.URL)
	fmt.Fprintln(w, "hello from roverd!")
}

func handleMove(w http.ResponseWriter, req *http.Request) {
	log.Println("new move request")

	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("failed to read request body:", err)
	}

	log.Printf("request body: %s\n", b)

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Fatalln("failed to unmarshal HTTP request body into map[string]interface{}:", err)
	}

	err = move.HandleMove(m)
	if err != nil {
		log.Println("failed to handle move:", err)
		w.WriteHeader(http.StatusInternalServerError)

		body := make(map[string]interface{})
		body["message"] = err.Error()

		b, err := json.Marshal(body)
		if err != nil {
			log.Fatalln("failed to marshal json error response:", err)
		}

		_, err = w.Write(b)
		if err != nil {
			log.Fatalln("failed to write json error response body:", err)
		}

		return
	}

	body := make(map[string]interface{})
	body["message"] = "success"

	b, err = json.Marshal(body)
	if err != nil {
		log.Fatalln("failed to marshal json response:", err)
	}

	_, err = w.Write(b)
	if err != nil {
		log.Fatalln("failed to write json response body:", err)
	}
}
