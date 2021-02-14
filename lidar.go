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
// Returns a pointer to exec.Cmd, which can be used to
// do something with lidar-scan process.
func StartLidar() (*exec.Cmd, error) {
	lidarChan = make(chan struct{}, 0)

	command := exec.Command("lidar-scan", "/dev/ttyUSB0")

	cmdStdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatalln(err, "failed to connect to lidar-scan's stdout")
	}

	cmdStderr, err := command.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to lidar-scan's stderr")
	}

	logFile, err := os.Create("lidar.log")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lidar-scan log file")
	}

	go scan(logFile, bufio.NewScanner(cmdStdout))
	go scan(logFile, bufio.NewScanner(cmdStderr))

	err = command.Start()
	fmt.Println("PROCESS PID:", command.Process.Pid)

	pidFile, err := os.Create("lidar-scan.pid")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create lidar-scan pid file")
	}
	pidFile.WriteString(strconv.Itoa(command.Process.Pid))

	if err != nil {
		return nil, errors.Wrap(err, "master: failed to start lidar-scan")
	}

	return command, nil
}

// StopLidar sends SIGINT to lidar-scan (if running).
func StopLidar() error {

	pidFile, err := os.Open("lidar-scan.pid")
	if err != nil {
		return errors.Wrap(err, "failed to open lidar-scan pid file")
	}

	pidStr, err := ioutil.ReadAll(pidFile)
	if err != nil {
		return errors.Wrap(err, "failed to read from lidar-scan pid file")
	}

	pid, err := strconv.Atoi(string(pidStr))
	if err != nil {
		return errors.Wrap(err, "failed to convert lidar-scan's pid to int")
	}

	syscall.Kill(pid, syscall.SIGINT)
	return nil
}

func scan(w io.Writer, scanner *bufio.Scanner) {
	for scanner.Scan() {
		t := time.Now().UnixNano()
		fmt.Fprintf(w, "%d: %s\n", t, scanner.Text())
	}
}
