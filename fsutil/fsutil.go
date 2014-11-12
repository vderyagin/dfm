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

		if relPath, _ := filepath.Rel(dir, path); strings.HasPrefix(relPath, ".") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

// IsEmptyDir determines whether given path corresponds to an empty directory.
func IsEmptyDir(path string) bool {
	if stat, err := os.Stat(path); !(err == nil && stat.IsDir()) {
		return false
	}

	entries, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return false
	}

	return entries == nil
}

// DeleteEmptyDirs remove empty directories starting at given path and going
// up filesystem hierarchy until it encounters a non-empty directory.
func DeleteEmptyDirs(start string) error {
	for dir := start; IsEmptyDir(dir); dir = filepath.Dir(dir) {
		if err := os.Remove(dir); err != nil {
			return err
		}
	}

	return nil
}

// Exists determines whether given path corresponds to an existing file.
func Exists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// IsRegularFile determines whether given path corresponds to a regular file.
func IsRegularFile(path string) bool {
	fi, err := os.Lstat(path)

	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// IsSymlink determines whether given path corresponds to a symbolic link.
func IsSymlink(path string) bool {
	if fi, err := os.Lstat(path); err != nil {
		return false
	} else if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}
