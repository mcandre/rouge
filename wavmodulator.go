package rouge

import (
	"github.com/go-audio/wav"

	"os"
)

// WavModulator passes out WAV file data.
type WavModulator struct {
	w *wav.Encoder
}

// NewWavModulator constructs a WavModulator.
func NewWavModulator(f *os.File, sampleRate int, bitDepth int, numChans int, audioFormat int) *WavModulator {
	return &WavModulator{w: wav.NewEncoder(f, sampleRate, bitDepth, numChans, audioFormat)}
}

// Encoder returns signal writers.
func (o *WavModulator) Encoder() (chan<- Message, <-chan error) {
	ch := make(chan Message)
	chErr := make(chan error)

	go func() {
		var m Message

		for {
			m = <-ch

			if m.Error != nil {
				break
			}

			//
			// DPSK...
			//

			// _, err := o.f.Write(m.Data)

			// if err != nil {
			// 	chErr<-err
			// 	break
			// }

			if m.Done {
				break
			}
		}
	}()

	return ch, chErr
}
