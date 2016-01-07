package dotfile

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
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

// String returns a string representation of State.
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
		formatStr = ansi.Color(" %s ", "green+b")
	case NotLinked:
		formatStr = ansi.Color(" %s ", "yellow+b")
	case Conflict:
		formatStr = ansi.Color(" %s ", "red+b")
	case Missing:
		formatStr = ansi.Color(" %s ", "red+bi")
	}

	return fmt.Sprintf(formatStr, *s)
}
