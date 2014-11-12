package logger

import "fmt"

// Logger logs performed DotFile actions
type Logger string

// New creates a pointer to a Logger object with a given id.
func New(id string) *Logger {
	l := Logger(id)
	return &l
}

// Success logs successful action
func (l *Logger) Success(msg string) {
	fmt.Printf("\x1b[1;32m%s\x1b[0m: %s\n", msg, l)
}

// Fail logs successful action
func (l *Logger) Fail(msg, reason string) {
	fmt.Printf("\x1b[1;31%s\x1b[0m: %s\n\t(%s)\n", msg, l, reason)
}

// Skip logs skipped action
func (l *Logger) Skip(msg, reason string) {
	fmt.Printf("\x1b[1;33%s\x1b[0m: %s\n\t(%s)\n", msg, l, reason)
}

// String representation of logger
func (l *Logger) String() string {
	return string(*l)
}
