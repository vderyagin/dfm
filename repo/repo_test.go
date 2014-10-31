package repo_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/repo"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	Describe("New", func() {
		It("Creates Repo object with absolute paths", func() {
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

		It("Returns empty collection for empty storage", func() {
			os.Mkdir("empty", 0777)
			absPath, _ := filepath.Abs("empty")

			repo := New(absPath, ".")

			Expect(repo.StoredDotFiles()).To(BeEmpty())
		})

		It("Returns collection including dotfile objects", func() {
			ioutil.WriteFile("bashrc", []byte{}, 0777)
			dfPath, _ := filepath.Abs("bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := repo.StoredDotFiles()

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(dfPath))
		})

		It("Returns collection including dotfile objects from nested directories", func() {
			os.MkdirAll("foo/bar/baz", 0777)
			ioutil.WriteFile("foo/bar/baz/bashrc", []byte{}, 0777)
			dfPath, _ := filepath.Abs("foo/bar/baz/bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := repo.StoredDotFiles()

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(dfPath))
		})
	})
})
