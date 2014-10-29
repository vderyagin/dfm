package dotfile

// DotFile type represents a single dotfile, defined by its storage and
// linking location. If dotfile is stored, its StoredLocation corresponds to a
// file within dotfiles repository and OriginalLocation - to symlink in user
// home directory where system expects original file to be.
type DotFile struct {
	OriginalLocation string
	StoredLocation   string
}
