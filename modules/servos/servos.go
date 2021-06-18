package servos

import (
	"fmt"
	"io"
	"log"

	"github.com/knei-knurow/frames"
)

const (
	TurnLeft  = "left"
	TurnRight = "right"
)

var (
	// Port is a serial port to which frames will be written.
	// It must be non-nil, otherwise this won't work. This means that
	// you should set this as soon as possible.
	Port io.ReadWriter

	frameHeader = [2]byte{'M', 'T'}
)

// TurnMove is a Move whose type is "turn".
type TurnMove struct {
	Degrees byte `json:"degrees"`
}

// ExecuteTurnMove packs a "turn move" into a frame and writes it to w.
func ExecuteTurnMove(move TurnMove) error {
	var (
		typeByte    byte = 'T'
		degreesByte byte = move.Degrees
	)

	var (
		frontSide byte = 'F'
		backSide  byte = 'B'
	)

	// we need to convert left and right (used by user)
	// to front 2 servos and back 2 servos (used by avr)
	frontDegrees := degreesByte
	backDegrees := 180 - degreesByte

	// write front servos frame
	data := []byte{typeByte, frontSide, frontDegrees}
	f := frames.Create(frameHeader, data)
	n, err := Port.Write(f)
	if err != nil {
		return fmt.Errorf("write frame to w: %v", err)
	}

	// TODO: add proper logging solution
	log.Printf("FRAME 1: wrote %d bytes to port\n", n)
	verbose := true
	if verbose {
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}

	// write back servos frame
	data = []byte{typeByte, backSide, backDegrees}
	f = frames.Create(frameHeader, data)
	n, err = Port.Write(f)
	if err != nil {
		return fmt.Errorf("write frame to w: %v", err)
	}

	log.Printf("FRAME 2: wrote %d bytes to port\n", n)
	if verbose {
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}

	return nil
}
