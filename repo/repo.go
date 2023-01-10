package repo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

// StoredDotFiles returns a channel producing DotFile objects for every stored
// dotfile.
func (r *Repo) StoredDotFiles() <-chan *dotfile.DotFile {
	dotFileChan := make(chan *dotfile.DotFile)

	go func(c chan<- *dotfile.DotFile) {
		for file := range fsutil.FilesIn(r.Store) {
			df := dotfile.DotFile{
				StoredLocation:   file,
				OriginalLocation: r.OriginalFilePath(file),
			}

			// is it generic file without conflicting host-specific one?
			noConflict := df.IsGeneric() && !fsutil.Exists(df.StoredLocation+host.DotFileLocalSuffix())

			if df.IsFromThisHost() || noConflict {
				c <- &df
			}
		}

		close(c)
	}(dotFileChan)

	return dotFileChan
}

// OriginalFilePath computes original path of dotfile (where it should be
// symlinked) based on path where it is stored.
func (r *Repo) OriginalFilePath(stored string) string {
	relPath, err := filepath.Rel(r.Store, stored)

	if err != nil {
		log.Fatal(err)
	}

	relPath = regexp.MustCompile(`\.force-copy`).ReplaceAllLiteralString(relPath, "")
	relPath = host.RemoveSuffix(relPath)

	return filepath.Join(r.Home, "."+relPath)
}

// StoredFilePath computes a path for stored dotfile corresponding to a given
// original path.
func (r *Repo) StoredFilePath(orig string, hostSpecific bool, forceCopy bool) (string, error) {
	relPath, err := filepath.Rel(r.Home, orig)

	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(relPath, ".") {
		return "", fmt.Errorf("%s is not a dotfile", orig)
	}

	// Handle case when file is host-local and already linked.
	if st, err := os.Readlink(orig); err == nil {
		if !filepath.IsAbs(st) {
			st = filepath.Join(filepath.Dir(orig), st)
		}

		if strings.HasSuffix(st, host.DotFileLocalSuffix()) {
			return st, nil
		}
	}

	storedRelPath := strings.TrimPrefix(relPath, ".")

	if hostSpecific {
		storedRelPath += host.DotFileLocalSuffix()
	}

	if forceCopy {
		storedRelPath += ".force-copy"
	}

	return filepath.Join(r.Store, storedRelPath), nil
}
