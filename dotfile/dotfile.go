package dotfile

import (
	"log"
	"os"
	"path/filepath"
)

// DotFile type represents a single dotfile, defined by its storage and
// linking location. If dotfile is stored, its StoredLocation corresponds to a
// file within dotfiles repository and OriginalLocation - to symlink in user
// home directory where system expects original file to be.
type DotFile struct {
	StoredLocation   string
	OriginalLocation string
}

// New returns a pointer to a DotFile object. Paths passed as arguments must
// be absolute.
func New(stored, original string) *DotFile {
	if !(filepath.IsAbs(stored) && filepath.IsAbs(original)) {
		log.Fatal("only absolute paths are accepted")
	}

	return &DotFile{
		StoredLocation:   stored,
		OriginalLocation: original,
	}
}

// IsStored returns true if given dotfile is properly stored and linked back
// to home dir, false otherwise.
func (df *DotFile) IsStored() bool {
	origInfo, err := os.Lstat(df.OriginalLocation)

	if err != nil {
		return false
	}

	origLinkTargetInfo, err := os.Stat(df.OriginalLocation)

	if err != nil {
		return false
	}

	storedInfo, err := os.Lstat(df.StoredLocation)

	if err != nil {
		return false
	}

	if origInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
		return false
	}

	if !storedInfo.Mode().IsRegular() {
		return false
	}

	return os.SameFile(origLinkTargetInfo, storedInfo)
}
