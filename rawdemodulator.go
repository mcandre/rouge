package rouge

import (
	"io"
	"os"
)

// RawDemodulator passes in raw file data.
type RawDemodulator struct {
	f *os.File
}

// NewRawDemodulator constructs a RawDemodulator.
func NewRawDemodulator(f *os.File) *RawDemodulator {
	return &RawDemodulator{f: f}
}

// Decoder returns signal readers.
func (o *RawDemodulator) Decoder() <-chan Message {
	ch := make(chan Message)

	go func() {
		for {
			buf := make([]byte, 1024)
			count, err := o.f.Read(buf)
			m := Message{ Data: buf[:count] }

			if err != nil {
				if err == io.EOF {
					m.Done = true
				} else {
					m.Error = &err
				}

				ch<-m
				break
			}

			ch<-m
		}
	}()

	return ch
}
