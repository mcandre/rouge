package rouge

// Demodulator reads signals.
type Demodulator interface {
	// Decode returns a channel for reading Messages.
	Decoder() <-chan Message

	// SampleRate is Hz.
	SampleRate() uint32

	// BitDepth, e.g. 16, 32.
	BitDepth() uint16

	// NumChannels, e.g. 1 (mono), 2 (stereo).
	NumChannels() uint16

	// WavCategory, e.g. 1 (PCM).
	WavCategory() uint16
}
