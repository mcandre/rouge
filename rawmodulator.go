package rouge

import (
	"log"
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
func (o *RawModulator) Encoder() (<-chan struct{}, chan<- Message, <-chan error) {
	chDone := make(chan struct{})
	ch := make(chan Message)
	chErr := make(chan error)

	go func() {
		defer func() {
			close(ch)
			close(chErr)

			if err := o.f.Close(); err != nil {
				log.Print(err)
			}
		}()

		var m Message

		for {
			m = <-ch

			if m.Error != nil {
				return
			}

			_, err := o.f.Write(m.Data)

			if err != nil {
				chErr<-err
				return
			}

			if m.Done {
				chDone<-struct{}{}
				return
			}
		}
	}()

	return chDone, ch, chErr
}
