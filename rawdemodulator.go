package rouge

import (
	"io"
	"log"
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
		defer func() {
			if err := o.f.Close(); err != nil {
				log.Print(err)
			}
		}()

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
				return
			}

			ch<-m
		}
	}()

	return ch
}

func (o RawDemodulator) SampleRate() uint32 {
	return 22050
}

func (o RawDemodulator) BitDepth() uint16 {
	return 16
}

func (o RawDemodulator) NumChannels() uint16 {
	return 1
}

func (o RawDemodulator) WavCategory() uint16 {
	return 1
}
