package testutil

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo"
)

// ExecuteEachInTempDir makes Ginkgo run each test in given context in
// temporary directory, which is cleaned up after each test run.
func ExecuteEachInTempDir() {
	var tempDir, startingDir string
	var err error

	if startingDir, err = os.Getwd(); err != nil {
		log.Fatalln("Failed to get current directory.")
	}

	ginkgo.BeforeEach(func() {
		var err error

		if tempDir, err = ioutil.TempDir(os.TempDir(), "dotfiletest"); err != nil {
			log.Fatalln("Failed to create temporary directory.")
		}

		if err := os.Chdir(tempDir); err != nil {
			log.Fatalf("Failed to cd into %s\n", tempDir)
		}
	})

	ginkgo.AfterEach(func() {
		os.Chdir(startingDir)
		os.RemoveAll(tempDir)
	})
}

// CreateDir creates a directory at a given path, also creating it's parent
// directories if needed. Loudly fails if anything goes wrong.
func CreateDir(path string) {
	os.MkdirAll(path, 0777)
}

// CreateFile creates an empty file at a given path, also creating it's parent
// directories if needed. Loudly fails if anything goes wrong.
func CreateFile(path string) {
	CreateDir(filepath.Dir(path))
	ioutil.WriteFile(path, []byte{}, 0777)
}
