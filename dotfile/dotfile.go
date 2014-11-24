package dotfile

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vderyagin/dfm/fsutil"
	"github.com/vderyagin/dfm/host"
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
	return fsutil.IsRegularFile(df.StoredLocation)
}

// IsLinked returns true if file is stored and linked to it's original
// location, false otherwise.
func (df *DotFile) IsLinked() bool {
	if !df.IsStored() {
		return false
	}

	if !fsutil.IsSymlink(df.OriginalLocation) {
		return false
	}

	origLinkTargetInfo, err := os.Stat(df.OriginalLocation)

	if err != nil {
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

// Store puts file in storage and links to it from original location.
func (df *DotFile) Store() error {
	if df.IsLinked() {
		return SkipError("is stored and linked already")
	}

	if !df.IsReadyToBeStored() {
		return FailError("can not be stored")
	}

	if err := os.MkdirAll(filepath.Dir(df.StoredLocation), 0777); err != nil {
		return FailError(err.Error())
	}

	if err := os.Rename(df.OriginalLocation, df.StoredLocation); err != nil {
		return FailError(err.Error())
	}

	if err := os.Symlink(df.StoredLocation, df.OriginalLocation); err != nil {
		return FailError(err.Error())
	}

	return nil
}

// Link links stored dotfile to its original location.
func (df *DotFile) Link() error {
	if !df.IsStored() {
		return errors.New("can link only already stored files")
	}

	if df.IsLinked() {
		return errors.New("is linked already")
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

// Restore moves stored file back into its original location, replacing symlink.
func (df *DotFile) Restore() error {
	if df.IsReadyToBeStored() {
		return SkipError("not stored to begin with")
	}

	if !df.IsLinked() {
		return FailError("can restore only properly linked files")
	}

	if err := os.Remove(df.OriginalLocation); err != nil {
		return FailError(err.Error())
	}

	if err := os.Rename(df.StoredLocation, df.OriginalLocation); err != nil {
		return FailError(err.Error())
	}

	if err := fsutil.DeleteEmptyDirs(filepath.Dir(df.StoredLocation)); err != nil {
		return FailError(err.Error())
	}

	return nil
}

// Delete removes stored file and link to it from home dir, fails if file is
// not linked properly.
func (df *DotFile) Delete() error {
	if !(fsutil.Exists(df.OriginalLocation) || fsutil.Exists(df.StoredLocation)) {
		return SkipError("does not exist")
	}

	if !df.IsLinked() {
		return FailError("can delete only properly linked files")
	}

	if err := os.Remove(df.StoredLocation); err != nil {
		return FailError(err.Error())
	}

	if err := os.Remove(df.OriginalLocation); err != nil {
		return FailError(err.Error())
	}

	if err := fsutil.DeleteEmptyDirs(filepath.Dir(df.StoredLocation)); err != nil {
		return FailError(err.Error())
	}

	if err := fsutil.DeleteEmptyDirs(filepath.Dir(df.OriginalLocation)); err != nil {
		return FailError(err.Error())
	}

	return nil
}

// IsFromThisHost returns true if dotfile is specific to current host, false
// otherwise.
func (df *DotFile) IsFromThisHost() bool {
	return strings.HasSuffix(df.StoredLocation, host.DotFileLocalSuffix())
}

// IsGeneric returns true if dotfile is not specific to any host, false
// otherwise.
func (df *DotFile) IsGeneric() bool {
	return !host.PathRegexp.MatchString(df.StoredLocation)
}
