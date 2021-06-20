package servos

import (
	"fmt"
	"log"
	"time"

	"github.com/knei-knurow/frames"
	"github.com/knei-knurow/roverd/sercom"
)

var (
	// Port is a serial port to which frames will be written. It must be non-nil.
	Port sercom.Serial

	frameHeader = [2]byte{'M', 'T'}
)

var Verbose bool

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

	// write servos frame for front wheels
	data1 := []byte{typeByte, frontSide, frontDegrees}
	f1 := frames.Create(frameHeader, data1)
	n, err := Port.WriteTimeout(f1, time.Second)
	if err != nil {
		return fmt.Errorf("write frame to port %v", err)
	}

	// TODO: add proper logging solution
	log.Printf("FRAME 1: wrote %d bytes to port\n", n)
	if Verbose {
		for _, b := range f1 {
			log.Println(frames.DescribeByte(b))
		}
	}

	// write servos frame for back wheels
	data2 := []byte{typeByte, backSide, backDegrees}
	f2 := frames.Create(frameHeader, data2)
	n, err = Port.WriteTimeout(f2, time.Second)
	if err != nil {
		return fmt.Errorf("write frame to port: %v", err)
	}

	log.Printf("FRAME 2: wrote %d bytes to port\n", n)
	if Verbose {
		for _, b := range f2 {
			log.Println(frames.DescribeByte(b))
		}
	}

	return nil
}
