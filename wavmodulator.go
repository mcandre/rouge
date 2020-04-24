package rouge

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"

	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

// WavModulator passes out WAV file data.
type WavModulator struct {
	w *wav.Encoder
	sourceSampleRate uint32
	sourceBitDepth uint16
	sourceNumChannels uint16
}

// NewWavModulator constructs a WavModulator.
func NewWavModulator(f *os.File, sourceSampleRate uint32, sourceBitDepth uint16, sourceNumChannels uint16, sampleRate uint32, bitDepth uint16, numChans uint16, audioFormat uint16) *WavModulator {
	return &WavModulator{
		w: wav.NewEncoder(f, int(sampleRate), int(bitDepth), int(numChans), int(audioFormat)),
		sourceSampleRate: sourceSampleRate,
		sourceBitDepth: sourceBitDepth,
		sourceNumChannels: sourceNumChannels,
	}
}

// Encoder returns signal writers.
func (o *WavModulator) Encoder() (chan<- Message, <-chan error) {
	ch := make(chan Message)
	chErr := make(chan error)

	go func() {
		defer func() {
			if err := o.w.Close(); err != nil {
				log.Print(err)
			}
		}()

		var m Message

		for {
			m = <-ch

			if m.Error != nil {
				return
			}

			//
			// DPSK...
			//

			var ds []int
			var d int64
			bb := bytes.NewBuffer(m.Data)

			for {
				err := binary.Read(bb, binary.LittleEndian, &d)

				if err != nil {
					if err == io.EOF {
						ds = append(ds, int(d))
						break
					} else {
						chErr<-err
						return
					}
				}

				ds = append(ds, int(d))
			}

			format := &audio.Format{
				NumChannels: int(o.sourceNumChannels),
				SampleRate: int(o.sourceSampleRate),
			}

			buf := audio.IntBuffer{
				Format: format,
				Data: ds,
				SourceBitDepth: int(o.sourceBitDepth),
			}

			if err := o.w.Write(&buf); err != nil {
				chErr<-err
				return
			}

			if m.Done {
				return
			}
		}
	}()

	return ch, chErr
}
