package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	command := exec.Command("lidar-scan", "/dev/ttyUSB0")

	cmdStdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatalln("master: failed to connect to command's stdout:", err)
	}

	cmdStderr, err := command.StderrPipe()
	if err != nil {
		log.Fatalln("master: failed to connect to command's stderr:", err)
	}

	go scan(bufio.NewScanner(cmdStdout))
	go scan(bufio.NewScanner(cmdStderr))

	err = command.Start()
	if err != nil {
		log.Fatalln("master: failed to run command:", err)
	}

	err = command.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "master: failed to wait for command", err)
		return
	}
}

func scan(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Printf("master: from command: %s\n", scanner.Text())
	}
}
