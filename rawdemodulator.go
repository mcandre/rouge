package rouge

import (
	"fmt"
	"io"
	"os"
)

// RawDemodulator passes in PCM samples from BO 32-bit unsigned integers.
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
		defer func() {
			close(ch)

			if err := o.f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		}()

		for {
			buf := make([]byte, 1024)
			count, err := o.f.Read(buf)
			m := Message{Data: buf[:count]}

			if err != nil {
				if err == io.EOF {
					m.Done = true
				} else {
					m.Error = &err
				}

				ch <- m
				return
			}

			ch <- m
		}
	}()

	return ch
}

// SampleRate unspecified.
func (o RawDemodulator) SampleRate() int {
	return 0
}

// BitDepth unspecified.
func (o RawDemodulator) BitDepth() int {
	return 0
}

// NumChannels unspecified.
func (o RawDemodulator) NumChannels() int {
	return 0
}

// WavCategory unspecified.
func (o RawDemodulator) WavCategory() int {
	return 0
}
