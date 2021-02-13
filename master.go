package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	command := exec.Command("lidar-scan", "/dev/ttyUSB0")

	out, err := command.CombinedOutput()
	if err != nil {
		log.Fatalln("master: failed to run command:", err)
	}

	fmt.Println("master: begin subprogram output")
	fmt.Println(string(out))
	fmt.Println("master: end subprogram output")

	// stdout, err := command.StdoutPipe()
	// if err != nil {
	// 	log.Fatalln("master: failed to connect to command's stdout:", err)
	// }

	// scanner := bufio.NewScanner(stdout)
	// go func() {
	// 	fmt.Println("in goroutine")
	// 	for scanner.Scan() {
	// 		fmt.Println("in scan")
	// 		fmt.Printf("\t > %s\n", scanner.Text())
	// 	}
	// }()

	// CUT

	// err = command.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "master: failed to wait for command", err)
	// 	return
	// }
}
