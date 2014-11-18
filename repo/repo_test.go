package repo_test

import (
	"os"
	"path/filepath"

	"github.com/vderyagin/dfm/dotfile"
	. "github.com/vderyagin/dfm/repo"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func chanToSlice(input <-chan *dotfile.DotFile) []*dotfile.DotFile {
	output := []*dotfile.DotFile{}

	for f := range input {
		output = append(output, f)
	}

	return output
}

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

			Expect(chanToSlice(repo.StoredDotFiles())).To(BeEmpty())
		})

		It("returns collection including dotfile objects", func() {
			CreateFile("bashrc")
			storedPath, _ := filepath.Abs("bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := chanToSlice(repo.StoredDotFiles())

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(storedPath))
		})

		It("returns collection including dotfile objects from nested directories", func() {
			CreateDir("foo/bar/baz")
			CreateFile("foo/bar/baz/bashrc")
			storedPath, _ := filepath.Abs("foo/bar/baz/bashrc")
			repoPath, _ := filepath.Abs(".")

			repo := New(repoPath, ".")
			dotfiles := chanToSlice(repo.StoredDotFiles())

			Expect(dotfiles).To(HaveLen(1))
			Expect(dotfiles[0].StoredLocation).To(Equal(storedPath))
		})

		It("returns multiple files, if there multiple stored", func() {
			repo := New(".", ".")
			CreateFile("foo")
			CreateFile("bar")

			Expect(chanToSlice(repo.StoredDotFiles())).To(HaveLen(2))
		})

		Context("host-specific dotfiles", func() {
			ExecuteEachWithHostName("myhost")

			It("returns only one of multiple files with same original location", func() {
				repo := New(".", ".")
				CreateFile("bashrc")
				CreateFile("bashrc.host-myhost")
				CreateFile("bashrc.host-otherhost")

				Expect(chanToSlice(repo.StoredDotFiles())).To(HaveLen(1))
			})

			It("ignores all dotfiles specific to other hosts", func() {
				repo := New(".", ".")
				CreateFile("bashrc.host-otherhost")
				CreateFile("bashrc.host-stillotherhost")

				Expect(chanToSlice(repo.StoredDotFiles())).To(BeEmpty())
			})

			It("favors files specific to current host over files from other host", func() {
				repo := New(".", ".")
				CreateFile("bashrc.host-myhost")
				CreateFile("bashrc.host-otherhost")

				expected, _ := filepath.Abs("bashrc.host-myhost")
				dotfiles := chanToSlice(repo.StoredDotFiles())
				Expect(dotfiles).To(HaveLen(1))
				Expect(dotfiles[0].StoredLocation).To(Equal(expected))
			})

			It("favors files specific to current host over generic ones", func() {
				repo := New(".", ".")
				CreateFile("bashrc")
				CreateFile("bashrc.host-myhost")

				expected, _ := filepath.Abs("bashrc.host-myhost")
				dotfiles := chanToSlice(repo.StoredDotFiles())
				Expect(dotfiles).To(HaveLen(1))
				Expect(dotfiles[0].StoredLocation).To(Equal(expected))
			})

			It("favors generic files over files specific to other hosts", func() {
				repo := New(".", ".")
				CreateFile("bashrc")
				CreateFile("bashrc.host-otherhost")

				expected, _ := filepath.Abs("bashrc")
				dotfiles := chanToSlice(repo.StoredDotFiles())
				Expect(dotfiles).To(HaveLen(1))
				Expect(dotfiles[0].StoredLocation).To(Equal(expected))
			})
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

		Context("host-specific dotfiles", func() {
			ExecuteEachWithHostName("myhost")

			It("removes host-specific suffix for current host", func() {
				orig := repo.OriginalFilePath(filepath.Join(repo.Store, "bashrc.host-myhost"))

				Expect(orig).To(Equal(filepath.Join(repo.Home, ".bashrc")))
			})

			It("removes host-specific suffix for other hosts", func() {
				orig := repo.OriginalFilePath(filepath.Join(repo.Store, "bashrc.host-otherhost"))

				Expect(orig).To(Equal(filepath.Join(repo.Home, ".bashrc")))
			})
		})
	})

	Describe("StoredFilePath", func() {
		repo := New("/store", "/")

		It("returns proper file name for simple case", func() {
			stored, err := repo.StoredFilePath(filepath.Join(repo.Home, ".bashrc"), false)

			Expect(err).To(Succeed())
			Expect(stored).To(Equal(filepath.Join(repo.Store, "bashrc")))
		})

		It("returns proper file name for for deeply nested file", func() {
			orig := filepath.Join(repo.Home, ".config/camlistore/server-config.json")
			stored, err := repo.StoredFilePath(orig, false)

			Expect(err).To(Succeed())
			Expect(stored).To(Equal(filepath.Join(repo.Store, "config/camlistore/server-config.json")))
		})

		It("fails if path from home directory does not start with dot", func() {
			df, err := repo.StoredFilePath(filepath.Join(repo.Home, "bashrc"), false)

			Expect(df).To(BeEmpty())
			Expect(err).NotTo(Succeed())
		})

		Context("host-specific dotfiles", func() {
			ExecuteEachWithHostName("myhost")

			It("returns name with host-specific suffix when requested", func() {
				df, err := repo.StoredFilePath(filepath.Join(repo.Home, ".bashrc"), true)

				Expect(err).To(Succeed())
				Expect(df).To(HaveSuffix(".host-myhost"))
			})

			Context("original file is a link", func() {
				ExecuteEachInTempDir()

				It("returns file path specific to current host, if linked", func() {
					repo := New(".", ".")
					CreateFile("foo.host-myhost")
					os.Symlink("foo.host-myhost", ".foo")

					stored, err := repo.StoredFilePath(filepath.Join(repo.Home, ".foo"), false)

					Expect(err).To(Succeed())
					Expect(stored).To(HaveSuffix(".host-myhost"))
				})

				It("returns generic file path if linked to file specific to other host", func() {
					repo := New(".", ".")
					CreateFile("foo.host-otherhost")
					os.Symlink("foo.host-otherhost", ".foo")

					stored, err := repo.StoredFilePath(filepath.Join(repo.Home, ".foo"), false)

					Expect(err).To(Succeed())
					Expect(stored).NotTo(HaveSuffix(".host-myhost"))
					Expect(stored).NotTo(HaveSuffix(".host-otherhost"))
					Expect(stored).To(HaveSuffix("/foo"))
				})
			})
		})
	})
})
