package fsutil

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FilesIn returns a collection of absolute paths of files in dir. Returns
// empty collection if provided argument is not a valid directory.
func FilesIn(dir string) []string {
	var files []string

	dir, err := filepath.Abs(dir)

	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		if basename := filepath.Base(path); strings.HasPrefix(basename, ".") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}
