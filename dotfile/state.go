package dotfile

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

// States DotFile can be in
var (
	Linked    = State("linked")
	NotLinked = State("not linked")
	Conflict  = State("conflict")
	Missing   = State("missing")
)

// State represents current state of DotFile
type State string

// String regurns strig representation o State
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
	var style string

	switch *s {
	case Linked:
		style = "green+b"
	case NotLinked:
		style = "yellow+b"
	case Conflict:
		style = "red+b"
	case Missing:
		style = "red+bi"
	}

	return ansi.Color(fmt.Sprintf(" %s ", *s), style)
}
