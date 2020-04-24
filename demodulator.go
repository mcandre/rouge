package rouge

// Demodulator reads signals.
type Demodulator interface {
	// Decode returns a channel for reading Messages.
	Decoder() <-chan Message
}
