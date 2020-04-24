package rouge

// Message encapsulates information.
type Message struct {
	// Error indicates an I/O problem.
	Error *error

	// Done indicates the end of the signal.
	Done bool

	// Data is all or part of the content.
	Data []byte
}
