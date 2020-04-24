package rouge

// Demodulator reads signals.
type Demodulator interface {
	// Decode returns a channel for reading Messages.
	Decoder() <-chan Message

	// SampleRate is Hz.
	SampleRate() int

	// BitDepth, e.g. 16, 32.
	BitDepth() int

	// NumChannels, e.g. 1 (mono), 2 (stereo).
	NumChannels() int

	// WavCategory, e.g. 1 (PCM).
	WavCategory() int
}
