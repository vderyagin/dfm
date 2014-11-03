package dotfile

import (
	"log"
	"path/filepath"
)

// DotFile type represents a single dotfile, defined by its storage and
// linking location. If dotfile is stored, its StoredLocation corresponds to a
// file within dotfiles repository and OriginalLocation - to symlink in user
// home directory where system expects original file to be.
type DotFile struct {
	OriginalLocation string
	StoredLocation   string
}

// New returns a pointer to a DotFile object. Paths passed as arguments must
// be absolute.
func New(original, stored string) *DotFile {

	if !(filepath.IsAbs(original) && filepath.IsAbs(stored)) {
		log.Fatal("only absolute paths are accepted")
	}

	return &DotFile{
		OriginalLocation: original,
		StoredLocation:   stored,
	}
}
