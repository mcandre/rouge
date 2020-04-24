package rouge

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/riff"
	"github.com/go-audio/wav"

	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// WavDemodulator passes in WAV file data.
type WavDemodulator struct {
	w *wav.Decoder
	p *riff.Parser
}

// NewWavDemodulator constructs a WavDemodulator.
func NewWavDemodulator(f *os.File) (*WavDemodulator, error) {
	w := wav.NewDecoder(f)

	if ok := w.IsValidFile(); !ok {
		return nil, fmt.Errorf("not a valid wav file")
	}

	p := riff.New(f)

	if err := w.FwdToPCM(); err != nil {
		return nil, err
	}

	return &WavDemodulator{w: w, p: p}, nil
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
				return
			}

			fmt.Fprintf(os.Stderr, "Next chunk...\n")

			buf := &audio.IntBuffer{
				Format: o.w.Format(),
				Data: make([]int, 1024),
			}

			count, err := o.w.PCMBuffer(buf)

			if err != nil && err != io.EOF {
				m.Error = &err
				ch<-m
				return
			}

			if count == 0 {
				m.Done = true
				ch<-m
				return
			}

			//
			// DPSK...
			//

			bb := new(bytes.Buffer)

			for _, d := range buf.Data {
				if err2 := binary.Write(bb, binary.LittleEndian, int64(d)); err2 != nil {
					m.Error = &err2
					ch<-m
					return
				}
			}

			m.Data = bb.Bytes()

			if err == io.EOF {
				m.Done = true
				ch<-m
				return
			}

			ch<-m
		}
	}()

	return ch
}

func (o WavDemodulator) SampleRate() uint32 {
	return o.w.SampleRate
}

func (o WavDemodulator) BitDepth() uint16 {
	return o.w.BitDepth
}

func (o WavDemodulator) NumChannels() uint16 {
	return o.w.NumChans
}

func (o WavDemodulator) WavCategory() uint16 {
	return o.w.WavAudioFormat
}
