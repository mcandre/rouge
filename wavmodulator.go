package rouge

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"

	"fmt"
	"os"
)

// WavModulatorConfig parameterizes a WavModulator.
type WavModulatorConfig struct {
	File        *os.File
	SampleRate  uint32
	BitDepth    uint16
	NumChannels uint16
	WavCategory uint16
}

// WavModulator passes out WAV file data.
type WavModulator struct {
	w      *wav.Encoder
	config WavModulatorConfig
}

// NewWavModulator constructs a WavModulator.
func NewWavModulator(config WavModulatorConfig) *WavModulator {
	return &WavModulator{
		w:      wav.NewEncoder(config.File, int(config.SampleRate), int(config.BitDepth), int(config.NumChannels), int(config.WavCategory)),
		config: config,
	}
}

// Encoder returns signal writers.
func (o *WavModulator) Encoder() (<-chan struct{}, chan<- Message, <-chan error) {
	chDone := make(chan struct{})
	ch := make(chan Message)
	chErr := make(chan error)

	go func() {
		defer func() {
			close(ch)
			close(chErr)

			if err := o.w.Close(); err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		}()

		var m Message

		for {
			m = <-ch

			if m.Error != nil {
				return
			}

			format := &audio.Format{
				NumChannels: int(o.config.NumChannels),
				SampleRate:  int(o.config.SampleRate),
			}

			buf := audio.IntBuffer{
				Format:         format,
				SourceBitDepth: int(o.config.BitDepth),
			}

			//
			// BPSK...
			//

			ys, err := BytesToUint32s(m.Data)

			if err != nil {
				chErr <- err
				return
			}

			var xs []int

			for _, y := range ys {
				xs = append(xs, int(y))
			}

			buf.Data = xs

			if err := o.w.Write(&buf); err != nil {
				chErr <- err
				return
			}

			if m.Done {
				chDone <- struct{}{}
				return
			}
		}
	}()

	return chDone, ch, chErr
}
