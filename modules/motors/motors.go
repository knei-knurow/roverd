package motors

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

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
	n, err := Port.Write(f)
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
	_, err = ReadTimeout(Port, res, time.Second)
	if err != nil {

	}

	return nil
}

// ReadTimeout reads from r into buf (just like io.Read).
// If the read operation takes more than timeout, it returns a non-nil error.
func ReadTimeout(r io.Reader, buf []byte, timeout time.Duration) (n int, err error) {
	type response struct {
		n   int
		err error
	}

	ticker := time.NewTicker(timeout)
	c := make(chan response)

	go func() {
		n, err := r.Read(buf)
		c <- response{n, err}
	}()

	select {
	case res := <-c:
		return res.n, res.err
	case <-ticker.C:
		return 0, errors.New("read timeout")
	}
}
