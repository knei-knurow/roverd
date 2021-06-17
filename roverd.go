package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/knei-knurow/roverd/handlers/move"
	"github.com/knei-knurow/roverd/modules/motors"
	"github.com/tarm/serial"
)

var (
	verbose bool
)

var (
	// Port on which the roverd will listen for requests
	port string

	// Port with the device controlling motors
	movePort io.ReadWriteCloser
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("roverd: ")

	flag.BoolVar(&verbose, "verbose", false, "print verbose output")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("failed to load .env file:", err)
	}

	port = os.Getenv("LISTEN_PORT")

	movePortName := os.Getenv("MOVE_PORT_NAME")
	baudRate, err := strconv.Atoi(os.Getenv("MOVE_BAUD_RATE"))
	if err != nil {
		log.Fatalf("cannot read baud rate: %v\n", err)
	}

	config := &serial.Config{
		Name: movePortName,
		Baud: baudRate,
	}

	movePort, err = serial.OpenPort(config)
	if err != nil {
		log.Fatalf("cannot open port %s: %v\n", movePortName, err)
	}

	motors.Port = movePort
}

func main() {
	fmt.Println("started")

	serveHTTP()
}

func serveHTTP() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/move", handleMove)
	http.ListenAndServe(":"+port, nil)
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	log.Printf("request from %s: %s %q", req.RemoteAddr, req.Method, req.URL)
	fmt.Fprintln(w, "hello from roverd!")
}

func handleMove(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("failed to read body:", err)
	}

	log.Printf("request body: %s\n", b)

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Fatalln("failed to unmarshal HTTP request body into map[string]interface{}:", err)
	}

	err = move.HandleMove(m)
	if err != nil {
		log.Fatalln("failed to handle move:", err)
	}
}
