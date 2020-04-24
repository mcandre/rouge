package rouge

import (
	"os"
)

// RawModulator passes out raw file data.
type RawModulator struct {
	f *os.File
}

// NewRawModulator constructs a RawModulator.
func NewRawModulator(f *os.File) *RawModulator {
	return &RawModulator{f: f}
}

// Encoder returns signal writers.
func (o *RawModulator) Encoder() (chan<- Message, <-chan error) {
	ch := make(chan Message)
	chErr := make(chan error)

	go func() {
		var m Message

		for {
			m = <-ch

			if m.Error != nil {
				break
			}

			_, err := o.f.Write(m.Data)

			if err != nil {
				chErr<-err
				break
			}

			if m.Done {
				break
			}
		}
	}()

	return ch, chErr
}
