package rouge

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

// ErrorInvalidQPSKSlice indicates a signal error.
var ErrorInvalidQPSKSlice = errors.New("Invalid QPSK slice (expected a sine or negated sine wave in <bit window> amplitude points)")

// QPSKDemodulator passes in mono PCM samples from BO 32-bit unsigned signed integers,
// reads 16-bit signed integers with two's complement,
// fits onto a quadrature phase shift signal,
// and finally outputs bit data in byte chunks,
// expanding each byte into a BO 32-bit unsigned integer.
type QPSKDemodulator struct {
	f            *os.File
	symbolWindow int
	bitBuffer    byte
	bitBufferLen uint8
	stemBuffer   []float64
	peakBuffer   []float64
	cursor       Cursor
}

// NewQPSKDemodulator constructs a QPSKDemodulator,
// given a peak amplitude threshold.
func NewQPSKDemodulator(f *os.File, symbolWindow int, innerThreshold float64, outerThreshold float64) *QPSKDemodulator {
	return &QPSKDemodulator{
		f:            f,
		symbolWindow: symbolWindow,
		cursor: Cursor{
			InnerThreshold: innerThreshold,
			OuterThreshold: outerThreshold,
		},
	}
}

// fitSignal attempts to read a bit from a QPSK sample slice.
// Returns an error on
func (o QPSKDemodulator) fitSignal() (*byte, error) {
	sign := math.Signbit(o.peakBuffer[0])

	var bitCandidate byte

	if !sign {
		bitCandidate = 1
	}

	for i := 1; i < len(o.peakBuffer); i++ {
		if math.Signbit(o.peakBuffer[i]) == sign {
			fmt.Fprintf(os.Stderr, "Cursor: %v\n", o.cursor)
			fmt.Fprintf(os.Stderr, "Peak buffer: %v\n", o.peakBuffer)
			return nil, ErrorInvalidQPSKSlice
		}

		sign = !sign
	}

	return &bitCandidate, nil
}

// Decoder returns signal readers.
func (o *QPSKDemodulator) Decoder() <-chan Message {
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

			var m Message

			if err != nil && err != io.EOF {
				m.Error = &err
				ch <- m
				return
			}

			m.Done = err == io.EOF

			if m.Done && count == 0 {
				ch <- m
				return
			}

			buf = buf[:count]

			if count%4 != 0 {
				m.Error = &ErrorInvalidBufferSize
				ch <- m
				return
			}

			for i := 0; i < count/4; i++ {
				sampleBytes := buf[i*4 : i*4+4]
				sample32s, err2 := BytesToUint32s(sampleBytes)

				if err2 != nil {
					m.Error = &err2
					ch <- m
					return
				}

				sampleTwosComplement := sample32s[0]

				var sample float64

				if sampleTwosComplement&0x80000000 == 0 {
					sample = float64(sampleTwosComplement)
				} else {
					sample = float64(int32(sampleTwosComplement))
				}

				o.cursor.Progress()

				if o.cursor.CheckStem(sample) {
					o.stemBuffer = append(o.stemBuffer, sample)
				}

				if !o.cursor.CheckPeak(sample) {
					continue
				}

				o.peakBuffer = append(o.peakBuffer, sample)

				// ...

				if len(o.peakBuffer) == o.symbolWindow {
					bitP, err2 := o.fitSignal()

					if err2 != nil {
						m.Error = &err2
						ch <- m
						return
					}

					bit := *bitP

					o.bitBuffer = (o.bitBuffer << 1) + bit
					o.bitBufferLen++

					if o.bitBufferLen == 8 {
						m.Data = []byte{o.bitBuffer}
						ch <- m

						if m.Done {
							return
						}

						o.bitBuffer = 0
						o.bitBufferLen = 0
					} else if m.Done {
						m.Error = &ErrorInvalidBitLength
						ch <- m
						return
					}

					o.peakBuffer = nil
				}
			}
		}
	}()

	return ch
}

// SampleRate unspecified.
func (o QPSKDemodulator) SampleRate() int {
	return 0
}

// BitDepth unspecified.
func (o QPSKDemodulator) BitDepth() int {
	return 0
}

// NumChannels unspecified.
func (o QPSKDemodulator) NumChannels() int {
	return 0
}

// WavCategory unspecified.
func (o QPSKDemodulator) WavCategory() int {
	return 0
}
