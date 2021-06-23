package main

import (
	"fmt"
	"os"
)

// Env represents all config vars from the .env file.
type Env struct {
	listenHost   string
	listenPort   string
	movePort     string
	movePortBaud string
}

func (e Env) String() string {
	return fmt.Sprint("listenHost:", e.listenHost, " listenPort:", e.listenPort, " movePort:", e.movePort, " movePortBaud:", e.movePortBaud)
}

// Load loads env vars from shell to e.
func (e *Env) Load() {
	e.listenHost = os.Getenv("ROVERD_LISTEN_HOST")
	e.listenPort = os.Getenv("ROVERD_LISTEN_PORT")
	e.movePort = os.Getenv("ROVERD_MOVE_PORT")
	e.movePortBaud = os.Getenv("ROVERD_MOVE_PORT_BAUD")
}
