package motors

import (
	"fmt"
	"io"
	"log"

	"github.com/knei-knurow/frames"
)

const (
	GoForward  = "forward"
	GoBackward = "backward"
	GoStop     = "stop"
)

var (
	// Port is a serial port to which frames will be written.
	// It must be non-nil, otherwise this won't work. This means that
	// you should set this as soon as possible.
	Port io.ReadWriter

	frameHeader = [2]byte{'M', 'T'}
)

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
	n, err := Port.Write(f)
	if err != nil {
		return fmt.Errorf("write frame to w: %v", err)
	}

	// TODO: add proper logging solution
	log.Printf("wrote %d bytes to port\n", n)
	verbose := true
	if verbose {
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}
	return nil
}
