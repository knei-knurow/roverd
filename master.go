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

	stdout, err := command.StderrPipe()
	if err != nil {
		log.Fatalln("master: failed to connect to command's stdout:", err)
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		fmt.Println("in goroutine")
		for scanner.Scan() {
			fmt.Println("in scan")
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	err = command.Start()
	if err != nil {
		log.Fatalln("maste: failed to run command:", err)
	}

	err = command.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "master: failed to wait for command", err)
		return
	}
}
