package repo_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/repo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
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

	It("Returns empty collection for empty storage", func() {
		os.Mkdir("empty", 0777)
		absPath, _ := filepath.Abs("empty")

		repo := Repo{Store: absPath}

		Expect(repo.StoredDotFiles()).To(BeEmpty())
	})

	It("Returns collection including dotfile objects", func() {
		ioutil.WriteFile("bashrc", []byte{}, 0777)
		dfPath, _ := filepath.Abs("bashrc")
		repoPath, _ := filepath.Abs(".")

		repo := Repo{Store: repoPath}
		dotfiles := repo.StoredDotFiles()

		Expect(dotfiles).To(HaveLen(1))
		Expect(dotfiles[0].StoredLocation).To(Equal(dfPath))
	})

	It("Returns collection including dotfile objects from nested directories", func() {
		os.MkdirAll("foo/bar/baz", 0777)
		ioutil.WriteFile("foo/bar/baz/bashrc", []byte{}, 0777)
		dfPath, _ := filepath.Abs("foo/bar/baz/bashrc")
		repoPath, _ := filepath.Abs(".")

		repo := Repo{Store: repoPath}
		dotfiles := repo.StoredDotFiles()

		Expect(dotfiles).To(HaveLen(1))
		Expect(dotfiles[0].StoredLocation).To(Equal(dfPath))
	})
})
