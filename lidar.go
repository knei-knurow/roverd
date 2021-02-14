package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// StartLidar runs lidar-scan.
func StartLidar() error {
	command := exec.Command("lidar-scan", "/dev/ttyUSB0")

	cmdStdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatalln(err, "failed to connect to lidar-scan's stdout")
	}

	cmdStderr, err := command.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to connect to lidar-scan's stderr")
	}

	logFile, err := os.Create("lidar.log")
	if err != nil {
		return errors.Wrap(err, "failed to create lidar log file")
	}

	go scan(logFile, bufio.NewScanner(cmdStdout))
	go scan(logFile, bufio.NewScanner(cmdStderr))

	fmt.Println("master: starting lidar-scan...")
	err = command.Start()
	if err != nil {
		return errors.Wrap(err, "master: failed to run lidar-scan")
	}

	err = command.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for lidar-scan to exit")
	}

	return nil
}

func scan(w io.Writer, scanner *bufio.Scanner) {
	for scanner.Scan() {
		t := time.Now().UnixNano()
		fmt.Fprintf(w, "%d: %s\n", t, scanner.Text())
	}
}
