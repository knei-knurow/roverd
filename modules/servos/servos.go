package servos

import (
	"fmt"
	"io"

	"github.com/knei-knurow/frames"
)

const (
	TurnLeft  = "left"
	TurnRight = "right"
)

// Port is a serial port to which frames will be written.
// It must be non-nil, otherwise this won't work.
var Port io.ReadWriter

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
