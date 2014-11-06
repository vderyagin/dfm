package repo_test

import (
	"path/filepath"

	. "github.com/vderyagin/dfm/repo"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	Describe("New", func() {
		It("creates Repo object with absolute paths", func() {
			relStore := "store"
			relHome := "home"
			absStore, _ := filepath.Abs(relStore)
			absHome, _ := filepath.Abs(relHome)

			repo := New(relStore, relHome)

			Expect(repo.Store).To(Equal(absStore))
			Expect(repo.Home).To(Equal(absHome))
		})
	})

	Describe("StoredDotFiles", func() {
		ExecuteEachInTempDir()

		It("returns empty collection for empty storage", func() {
			CreateDir("empty")
			absPath, _ := filepath.Abs("empty")

			repo := New(absPath, ".")

			Expect(repo.StoredDotFiles()).To(BeEmpty())
		})

		It("returns collection including dotfile objects", func() {
			CreateFile("bashrc")
			storedPath, _ := filepath.Abs("bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := repo.StoredDotFiles()

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(storedPath))
		})

		It("returns collection including dotfile objects from nested directories", func() {
			CreateDir("foo/bar/baz")
			CreateFile("foo/bar/baz/bashrc")
			storedPath, _ := filepath.Abs("foo/bar/baz/bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := repo.StoredDotFiles()

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(storedPath))
		})

		It("returns multiple files, if there multiple stored", func() {
			repo := New(".", ".")
			CreateFile("foo")
			CreateFile("bar")

			Expect(repo.StoredDotFiles()).To(HaveLen(2))
		})
	})

	Describe("OriginalFilePath", func() {
		repo := New("/store", "/")

		It("returns proper file name for simple case", func() {
			orig := repo.OriginalFilePath(filepath.Join(repo.Store, "bashrc"))
			Expect(orig).To(Equal(filepath.Join(repo.Home, ".bashrc")))
		})

		It("returns proper file name for for deeply nested file", func() {
			orig := repo.OriginalFilePath(filepath.Join(repo.Store, "config/camlistore/server-config.json"))
			Expect(orig).To(Equal(filepath.Join(repo.Home, ".config/camlistore/server-config.json")))
		})
	})
})
