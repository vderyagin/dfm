package dotfile

import (
	"fmt"
	"os"
)

// States DotFile can be in.
var (
	Linked    = State("linked")
	NotLinked = State("not linked")
	Conflict  = State("conflict")
	Missing   = State("missing")
)

// State represents current state of DotFile.
type State string

// String regurns string representation o State.
func (s State) String() string {
	return string(s)
}

// CurrentState returns State object representing a state DotFile is currently
// in.
func (df *DotFile) CurrentState() *State {
	if df.IsLinked() {
		return &Linked
	} else if !df.IsStored() {
		return &Missing
	} else if _, err := os.Lstat(df.OriginalLocation); os.IsNotExist(err) {
		return &NotLinked
	}

	return &Conflict
}

// ColorString generates ansi-colored representation of given State object.
func (s *State) ColorString() string {
	var formatStr string

	switch *s {
	case Linked:
		formatStr = "\x1b[1;32m %s \x1b[0m" // green bold
	case NotLinked:
		formatStr = "\x1b[1;33m %s \x1b[0m" // yellow bold
	case Conflict:
		formatStr = "\x1b[1;31m %s \x1b[0m" // red bold
	case Missing:
		formatStr = "\x1b[1;7;31m %s \x1b[0m" // red bold inverted
	}

	return fmt.Sprintf(formatStr, *s)
}
