package fsutil_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/fsutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fsutil", func() {
	var tempDir, startingDir string
	var err error

	if startingDir, err = os.Getwd(); err != nil {
		log.Fatalln("Failed to get current directory.")
	}

	BeforeEach(func() {
		var err error

		if tempDir, err = ioutil.TempDir(os.TempDir(), "dotfiletest"); err != nil {
			log.Fatalln("Failed to create temporary directory.")
		}

		if err := os.Chdir(tempDir); err != nil {
			log.Fatalf("Failed to cd into %s\n", tempDir)
		}
	})

	AfterEach(func() {
		os.Chdir(startingDir)
		os.RemoveAll(tempDir)
	})

	Describe("FilesIn", func() {
		It("Returns empty collection if argument does not exist", func() {
			Expect(FilesIn("nonexistent_dir")).To(BeEmpty())
		})

		It("Returns empty collection for empty directory", func() {
			Expect(FilesIn(".")).To(BeEmpty())
		})

		It("Returns empty collection for directory with hidden files", func() {
			ioutil.WriteFile(".hidden", []byte{}, 0777)

			Expect(FilesIn(".")).To(BeEmpty())
		})

		It("Returns collection including regular non-hidden files", func() {
			baseName := "bashrc"
			absPath, _ := filepath.Abs(baseName)
			ioutil.WriteFile(baseName, []byte{}, 0777)

			Expect(FilesIn(".")).To(ContainElement(absPath))
		})

		It("Returns collection including deeply nested files", func() {
			relPath := "config/camlistore/server-config.json"
			absPath, _ := filepath.Abs(relPath)
			os.MkdirAll(filepath.Dir(absPath), 0777)
			ioutil.WriteFile(relPath, []byte{}, 0777)

			Expect(FilesIn(".")).To(ContainElement(absPath))
		})
	})
})
