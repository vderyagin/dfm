package testutil

import (
	"io/ioutil"
	"log"
	"os"

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
