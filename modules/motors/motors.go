package motors

import (
	"fmt"
	"io"

	"github.com/knei-knurow/frames"
)

const (
	GoForward  = "forward"
	GoBackward = "backward"
	GoStop     = "stop"

	TurnLeft  = "left"
	TurnRight = "right"
)

// Port is a serial port to which frames will be written.
// It must be non-nil, otherwise this won't work.
var Port io.ReadWriter

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

// TurnMove is a Move whose type is "turn".
type TurnMove struct {
	Side    string `json:"side"`
	Degrees byte   `json:"degrees"`
}

// SideByte translates "side" field from JSON request body
// to byte. It returns 0 in case the field is empty.
func (m TurnMove) SideByte() (directionByte byte) {
	if m.Side == TurnLeft {
		directionByte = 'L'
	} else if m.Side == TurnRight {
		directionByte = 'R'
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
	f := frames.Create([2]byte{'M', 'T'}, data)
	_, err := Port.Write(f)
	if err != nil {
		return fmt.Errorf("write frame to w: %v", err)
	}

	// TODO: add proper logging solution
	// verbose := true
	// if verbose {
	// 	for _, b := range f {
	// 		log.Println(frames.DescribeByte(b))
	// 	}
	// }
	return nil
}

// ExecuteTurnMove packs a "turn move" into a frame and writes it to w.
func ExecuteTurnMove(move TurnMove) error {
	var (
		typeByte    byte = 'T'
		sideByte    byte = move.SideByte()
		degreesByte byte = move.Degrees
	)

	data := []byte{typeByte, sideByte, degreesByte}
	f := frames.Create([2]byte{'M', 'T'}, data)
	_, err := Port.Write(f)
	if err != nil {
		return fmt.Errorf("write frame to w: %v", err)
	}

	return nil
}
