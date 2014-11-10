package repo

import (
	"fmt"
	"log"
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

// StoredDotFiles returns a collection, containing DotFile object for each
// stored dotfile.
func (r *Repo) StoredDotFiles() []*dotfile.DotFile {
	var dotfiles []*dotfile.DotFile

	for _, file := range fsutil.FilesIn(r.Store) {
		df := dotfile.DotFile{
			StoredLocation:   file,
			OriginalLocation: r.OriginalFilePath(file),
		}

		dotfiles = append(dotfiles, &df)
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

	relPath = regexp.MustCompile(`\.host-[[:alpha:]]+\z`).ReplaceAllLiteralString(relPath, "")

	return filepath.Join(r.Home, "."+relPath)
}

// StoredFilePath computes a path for stored dotfile corresponding to a given
// original path.
func (r *Repo) StoredFilePath(orig string, hostSpecific bool) (string, error) {
	relPath, err := filepath.Rel(r.Home, orig)

	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(relPath, ".") {
		return "", fmt.Errorf("%s is not a dotfile", orig)
	}

	storedRelPath := strings.TrimPrefix(relPath, ".")

	if hostSpecific {
		storedRelPath += hostSpecificSuffix()
	}

	return filepath.Join(r.Store, storedRelPath), nil
}

func hostSpecificSuffix() string {
	return ".host-" + host.Name()
}
