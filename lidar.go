package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

var (
	lidarChan chan struct{}
)

// StartLidar runs lidar-scan.
func StartLidar() error {
	lidarChan = make(chan struct{}, 0)

	command := exec.Command("lidar-scan", "/dev/ttyUSB0")

	cmdStdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatalln(err, "failed to connect to lidar-scan's stdout")
	}

	cmdStderr, err := command.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to connect to lidar-scan's stderr")
	}

	logFile, err := os.Create("lidar-scan.log")
	if err != nil {
		return errors.Wrap(err, "failed to create lidar-scan log file")
	}

	go scan(logFile, bufio.NewScanner(cmdStdout))
	go scan(logFile, bufio.NewScanner(cmdStderr))

	err = command.Start()
	if err != nil {
		return errors.Wrapf(err, "failed to start lidar-scan")
	}
	log.Printf("started lidar-scan (pid %d)", command.Process.Pid)

	pidFile, err := os.Create("lidar-scan.pid")
	if err != nil {
		return errors.Wrap(err, "failed to create lidar-scan pid file")
	}
	defer pidFile.Close()

	n, err := fmt.Fprintf(pidFile, "%d", command.Process.Pid)
	if err != nil {
		return errors.Wrap(err, "failed to write lidar-scan pid to file")
	}
	log.Printf("wrote lidar-scan pid (%d) to file (%d bytes)", command.Process.Pid, n)

	return nil
}

// StopLidar sends SIGINT to lidar-scan (if running).
func StopLidar() error {
	pidFile, err := os.Open("lidar-scan.pid")
	if err != nil {
		return errors.Wrap(err, "failed to open lidar-scan pid file")
	}
	defer pidFile.Close()

	pidStr, err := ioutil.ReadAll(pidFile)
	if err != nil {
		return errors.Wrap(err, "failed to read from lidar-scan pid file")
	}

	pid, err := strconv.Atoi(string(pidStr))
	if err != nil {
		return errors.Wrap(err, "failed to convert lidar-scan's pid to int")
	}

	err = syscall.Kill(pid, syscall.SIGINT)
	if err != nil {
		return errors.Wrapf(err, "failed to send SIGINT to lidar-scan (pid %d)", pid)
	}

	log.Print("stopped lidar-scan")

	return nil
}

func scan(w io.Writer, scanner *bufio.Scanner) {
	// TODO: Add channel to signal end of scanning (sync and close files)
	for scanner.Scan() {
		t := time.Now().UnixNano()
		fmt.Fprintf(w, "%d: %s\n", t, scanner.Text())
	}
}
