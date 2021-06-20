package motors

import (
	"fmt"
	"log"
	"time"

	"github.com/knei-knurow/frames"
	"github.com/knei-knurow/roverd/ports"
)

const (
	GoForward  = "forward"
	GoBackward = "backward"
	GoStop     = "stop"
)

var (
	// Port is a serial port to which frames will be written. It must be non-nil.
	Port ports.Serial

	frameHeader = [2]byte{'M', 'T'}
)

var Verbose bool

// GoMove is a Move whose type is "go".
type GoMove struct {
	Direction string `json:"direction"`
	Speed     byte   `json:"speed"`
}

// DirectionByte translates "direction" field from JSON request body
// to byte. It returns 0 in case the field is empty.
func (m GoMove) DirectionByte() (directionByte byte) {
	if m.Direction == GoForward {
		directionByte = 'F'
	} else if m.Direction == GoBackward {
		directionByte = 'B'
	} else if m.Direction == GoStop {
		directionByte = 'S'
	}

	return
}

// ExecuteGoMove packs a "go move" into a frame and writes it to w.
func ExecuteGoMove(move GoMove) error {
	var (
		typeByte      byte = 'G'
		directionByte byte = move.DirectionByte()
		speedByte     byte = move.Speed
	)

	data := []byte{typeByte, directionByte, speedByte}
	f := frames.Create(frameHeader, data)
	n, err := Port.WriteTimeout(f, time.Second)
	if err != nil {
		return fmt.Errorf("write frame to port: %v", err)
	}

	// TODO: add proper logging solution
	// TODO: verify crc using package frames

	if Verbose {
		log.Printf("wrote %d bytes to port\n", n)
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}

	log.Println("waiting for response frame...")
	res := make([]byte, 8)
	_, err = Port.ReadTimeout(res, time.Second)
	if err != nil {
		return fmt.Errorf("read frame from port: %v", err)
	}

	return nil
}
