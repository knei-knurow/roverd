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
	"github.com/knei-knurow/frames"
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

	log.Printf("body: %s\n", b)

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Fatalln("failed to unmarshal:", err)
	}

	speed, ok := m["speed"].(float64)
	if !ok {
		log.Fatalln("failed to convert speed to float64")
	}

	data := []byte{'G', byte(speed)}
	f := frames.Create([2]byte{'M', 'T'}, data)
	_, err = movePort.Write(f)
	if err != nil {
		log.Fatalln("failed to write frame to port:", err)
	}

	if verbose {
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}
}
