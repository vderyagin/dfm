package dotfile

import (
	"errors"
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

// IsStored returns true if given dotfile is stored.
func (df *DotFile) IsStored() bool {
	if storedInfo, err := os.Lstat(df.StoredLocation); err != nil {
		return false
	} else if !storedInfo.Mode().IsRegular() {
		return false
	}

	return true
}

// IsLinked returns true if file is stored and linked to it's original
// location, false otherwise.
func (df *DotFile) IsLinked() bool {
	if !df.IsStored() {
		return false
	}

	origInfo, err := os.Lstat(df.OriginalLocation)

	if err != nil {
		return false
	}

	origLinkTargetInfo, err := os.Stat(df.OriginalLocation)

	if err != nil {
		return false
	}

	if origInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
		return false
	}

	storedInfo, _ := os.Lstat(df.StoredLocation)

	return os.SameFile(origLinkTargetInfo, storedInfo)
}

// IsReadyToBeStored returns true if dotfile is ready to be stored, that is if
// it is a regular file not conflicting with any of already stored files.
func (df *DotFile) IsReadyToBeStored() bool {
	origInfo, err := os.Lstat(df.OriginalLocation)

	if err != nil {
		return false
	}

	if !origInfo.Mode().IsRegular() {
		return false
	}

	if _, err := os.Lstat(df.StoredLocation); !os.IsNotExist(err) {
		return false
	}

	return true
}

// Store puts file in storage and links it back from there.
func (df *DotFile) Store() error {
	if !df.IsReadyToBeStored() {
		return errors.New("can not store")
	}

	if err := os.MkdirAll(filepath.Dir(df.StoredLocation), 0777); err != nil {
		return err
	}

	if err := os.Rename(df.OriginalLocation, df.StoredLocation); err != nil {
		return err
	}

	if err := os.Symlink(df.StoredLocation, df.OriginalLocation); err != nil {
		return err
	}

	return nil
}

// Link links stored dotfile to its original location.
func (df *DotFile) Link() error {
	if !df.IsStored() {
		return errors.New("file is not even stored")
	}

	if df.IsLinked() {
		return errors.New("file is linked already")
	}

	if _, err := os.Lstat(df.OriginalLocation); !os.IsNotExist(err) {
		return errors.New("conflicting file at original location")
	}

	if err := os.MkdirAll(filepath.Dir(df.OriginalLocation), 0777); err != nil {
		return err
	}

	if err := os.Symlink(df.StoredLocation, df.OriginalLocation); err != nil {
		return err
	}

	return nil
}
