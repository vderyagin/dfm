package repo

import (
	"log"
	"path/filepath"

	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/fsutil"
)

// Repo represents a place where dotfiles are stored.
type Repo struct{ Store, Home string }

// New returns a pointer to new instance of Repo. Makes sure that paths Repo
// initialized with are absolute, panics if they are not absolute and can not
// be made that.
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
		df := dotfile.DotFile{StoredLocation: file}
		dotfiles = append(dotfiles, &df)
	}

	return dotfiles
}
