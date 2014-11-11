package repo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/fsutil"
	"github.com/vderyagin/dfm/host"
)

// Repo represents a place where dotfiles are stored.
type Repo struct{ Store, Home string }

// New returns a pointer to new instance of Repo. Makes sure that paths Repo
// initialized with are absolute, fails loudly if they are not absolute and
// can not be made that.
func New(store, home string) *Repo {
	var absStore, absHome string
	var err error

	if absStore, err = filepath.Abs(store); err != nil {
		log.Fatal(err)
	}

	if absHome, err = filepath.Abs(home); err != nil {
		log.Fatal(err)
	}

	return &Repo{
		Store: absStore,
		Home:  absHome,
	}
}

// StoredDotFiles returns a collection, containing DotFile object for each
// stored dotfile.
func (r *Repo) StoredDotFiles() []*dotfile.DotFile {
	// Populate map indexed by dotfile original path to be able to weed out
	// conflicts - multiple dotfiles attempting to be linked from the same path
	// in home directory.
	dfMap := make(map[string]*dotfile.DotFile)

	for _, file := range fsutil.FilesIn(r.Store) {
		df := dotfile.DotFile{
			StoredLocation:   file,
			OriginalLocation: r.OriginalFilePath(file),
		}

		if _, clash := dfMap[df.OriginalLocation]; !clash && df.IsGeneric() || df.IsFromThisHost() {
			dfMap[df.OriginalLocation] = &df
		}
	}

	// Collect values from map into slice and return it.
	var dotfiles []*dotfile.DotFile

	for _, df := range dfMap {
		dotfiles = append(dotfiles, df)
	}

	return dotfiles
}

// OriginalFilePath computes original path of dotfile (where it should be
// symlinked) based on path where it is stored.
func (r *Repo) OriginalFilePath(stored string) string {
	relPath, err := filepath.Rel(r.Store, stored)

	if err != nil {
		log.Fatal(err)
	}

	relPath = host.RemoveSuffix(relPath)

	return filepath.Join(r.Home, "."+relPath)
}

// StoredFilePath computes a path for stored dotfile corresponding to a given
// original path.
func (r *Repo) StoredFilePath(orig string, hostSpecific bool) (string, error) {
	relPath, err := filepath.Rel(r.Home, orig)

	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(relPath, ".") {
		return "", fmt.Errorf("%s is not a dotfile", orig)
	}

	stat, err := os.Lstat(orig)

	if err == nil && stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		if st, err := os.Readlink(orig); err == nil {
			if !filepath.IsAbs(st) {
				st = filepath.Join(filepath.Dir(orig), st)
			}

			if strings.HasSuffix(st, host.DotFileSuffix()) {
				return st, nil
			}
		}
	}

	storedRelPath := strings.TrimPrefix(relPath, ".")

	if hostSpecific {
		storedRelPath += host.DotFileSuffix()
	}

	return filepath.Join(r.Store, storedRelPath), nil
}
