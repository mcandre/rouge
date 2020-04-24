package rouge

// Modulator writes signals.
type Modulator interface {
	// Encode returns a channel for writing Messages,
	// as well as channel for downstream I/O errors.
	Encoder() (chan<- Message, <-chan error)
}
