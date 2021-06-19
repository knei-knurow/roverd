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
		return fmt.Errorf("write frame to port: %v", err)
	}

	// TODO: add proper logging solution
	// TODO: verify crc using package frames

	verbose := true
	if verbose {
		log.Printf("wrote %d bytes to port\n", n)
		for _, b := range f {
			log.Println(frames.DescribeByte(b))
		}
	}

	// TODO: fix possible infinite waiting. Wait for max e.g 0.5s and then fail
	resData := make([]byte, 8)
	n, err = Port.Read(resData)
	if err != nil {
		return fmt.Errorf("read frame from port: %v", err)
	}

	if verbose {
		log.Printf("read %d bytes from port\n", n)
		for _, b := range resData {
			log.Println(frames.DescribeByte(b))
		}
	}

	return nil
}
