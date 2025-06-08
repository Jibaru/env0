package prompt

// Reader represents a minimal interface for reading input
type Reader interface {
	ReadString(delim byte) (string, error)
}
