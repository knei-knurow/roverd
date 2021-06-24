package ports

import (
	"errors"
	"io"
	"time"
)

// Serial represents a serial port. It implements the most basic uses of serial port.
type Serial interface {
	io.ReadWriteCloser // serial port abstraction
	ReadTimeout(buf []byte, timeout time.Duration) (n int, err error)
	WriteTimeout(buf []byte, timeout time.Duration) (n int, err error)
}

type SerialPort struct {
	io.ReadWriteCloser
}

type response struct {
	n   int
	err error
}

// ReadTimeout reads up to len(buf) bytes into buf (just like io.Read).
// If the read operation takes more than timeout, it returns a non-nil error.
func (s SerialPort) ReadTimeout(buf []byte, timeout time.Duration) (n int, err error) {
	ticker := time.NewTicker(timeout)
	c := make(chan response)

	go func() {
		n, err := s.Read(buf)
		c <- response{n, err}
	}()

	select {
	case res := <-c:
		return res.n, res.err
	case <-ticker.C:
		return 0, errors.New("read timeout")
	}
}

// WriteTimeout writes from buf to the underlying data stream (just like io.Write).
// If the write operation takes more than timeout, it returns a non-nil error.
func (s SerialPort) WriteTimeout(buf []byte, timeout time.Duration) (n int, err error) {
	ticker := time.NewTicker(timeout)
	c := make(chan response)

	go func() {
		n, err := s.Write(buf)
		c <- response{n, err}
	}()

	select {
	case res := <-c:
		return res.n, res.err
	case <-ticker.C:
		return 0, errors.New("write timeout")
	}
}
