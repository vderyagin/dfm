package repo

import (
	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/fsutil"
)

// Repo represents a place where dotfiles are stored.
type Repo struct {
	Store string
	Home  string
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
