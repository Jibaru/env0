package scripts

// Logger represents a minimal logging interface
type Logger interface {
	Printf(format string, v ...interface{})
}
