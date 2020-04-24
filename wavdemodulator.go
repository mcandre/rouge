package rouge

import (
	"github.com/go-audio/wav"

	"fmt"
	"os"
)

// WavDemodulator passes in WAV file data.
type WavDemodulator struct {
	w *wav.Decoder
}

// NewWavDemodulator constructs a WavDemodulator.
func NewWavDemodulator(f *os.File) (*WavDemodulator, error) {
	w := wav.NewDecoder(f)

	if ok := w.IsValidFile(); !ok {
		return nil, fmt.Errorf("not a valid wav file")
	}

	if err := w.FwdToPCM(); err != nil {
		return nil, err
	}

	return &WavDemodulator{w: w}, nil
}

// Decoder returns signal readers.
func (o *WavDemodulator) Decoder() <-chan Message {
	ch := make(chan Message)

	go func() {
		for {
			m := Message{}

			if o.w.EOF() {
				m.Done = true
				ch<-m
				break
			}

			chunk, err := o.w.NextChunk()



			fmt.Printf("Chunk: %v\n", chunk)


			//
			// DPSK...
			//

			if err != nil {
				// go-audio EOF identifier.
				if err.Error() == "error reading chunk header - EOF" {
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
