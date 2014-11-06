package dotfile_test

import (
	"log"
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/dotfile"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dotfile", func() {
	ExecuteEachInTempDir()

	stored := func() string {
		s, err := filepath.Abs("foo")
		if err != nil {
			log.Fatal(err)
		}
		return s
	}

	orig := func() string {
		o, err := filepath.Abs(".foo")
		if err != nil {
			log.Fatal(err)
		}
		return o
	}

	df := func() *DotFile {
		return New(stored(), orig())
	}

	Describe("IsStored", func() {
		It("returns true if file is properly stored", func() {
			CreateFile(stored())

			Expect(df().IsStored()).To(BeTrue())
		})

		It("returns false if stored file location is empty", func() {
			CreateFile(stored())
			os.Symlink(stored(), orig())
			os.Remove(stored())

			Expect(df().IsStored()).To(BeFalse())
		})

		It("returns false if stored file is not a regular file", func() {
			CreateDir(stored())
			os.Symlink(stored(), orig())

			Expect(df().IsStored()).To(BeFalse())
		})
	})

	Describe("IsLinked", func() {
		It("returns true if file is stored and linked properly", func() {
			CreateFile(stored())
			os.Symlink(stored(), orig())

			Expect(df().IsLinked()).To(BeTrue())
		})

		It("returns false if file is not stored properly", func() {
			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's nothing at original location", func() {
			CreateFile(stored())

			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's no symlink at original location", func() {
			CreateFile(stored())
			CreateDir(orig())

			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's wrong symlink at original location", func() {
			CreateFile(stored())
			CreateFile("wrong_file")
			os.Symlink("wrong_file", orig())

			Expect(df().IsLinked()).To(BeFalse())
		})
	})

	Describe("IsReadyToBeStored", func() {
		It("returns true if regular file not conflicting with stored ones", func() {
			CreateFile(orig())

			Expect(df().IsReadyToBeStored()).To(BeTrue())
		})

		It("returns false if original location is empty", func() {
			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})

		It("returns false if original location is not a regular file", func() {
			CreateDir(orig())

			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})

		It("returns false if conflicts with already stored file", func() {
			CreateFile(orig())
			CreateFile(stored())

			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})
	})

	Describe("Store", func() {
		It("stores file", func() {
			CreateFile(orig())

			Expect(df().IsStored()).To(BeFalse())
			Expect(df().Store()).To(Succeed())
			Expect(df().IsStored()).To(BeTrue())
		})

		It("creates intermediate directories for nested file", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			CreateDir(filepath.Dir(orig))
			CreateFile(orig)

			Expect(df.IsStored()).To(BeFalse())
			Expect(df.Store()).To(Succeed())
			Expect(df.IsStored()).To(BeTrue())
		})

		It("fails if file is not ready to be stored", func() {
			Expect(df().Store()).NotTo(Succeed())
			Expect(df().IsStored()).To(BeFalse())
		})
	})

	Describe("Link", func() {
		It("symlinks stored file to it's original location", func() {
			CreateFile(stored())

			Expect(df().IsLinked()).To(BeFalse())
			Expect(df().Link()).To(Succeed())
			Expect(df().IsLinked()).To(BeTrue())
		})

		It("creates deeply nested directories if necessary", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			CreateDir(filepath.Dir(stored))
			CreateFile(stored)

			Expect(df.IsLinked()).To(BeFalse())
			Expect(df.Link()).To(Succeed())
			Expect(df.IsLinked()).To(BeTrue())
		})

		It("fails if file is not stored", func() {
			Expect(df().Link()).NotTo(Succeed())
		})

		It("fails if file is already linked", func() {
			CreateFile(stored())
			os.Symlink(stored(), orig())

			Expect(df().Link()).NotTo(Succeed())
		})

		It("fails if there's conficting file at original location", func() {
			CreateFile(stored())
			CreateFile(orig())

			Expect(df().Link()).NotTo(Succeed())
		})
	})

	Describe("Restore", func() {
		It("restores the file in its original place", func() {
			CreateFile(stored())
			df().Link()

			Expect(df().IsLinked()).To(BeTrue())
			Expect(df().Restore()).To(Succeed())
			Expect(df().IsLinked()).To(BeFalse())
			Expect(df().IsReadyToBeStored()).To(BeTrue())
		})

		It("deletes any empty intermediate directories in store", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			CreateDir(filepath.Dir(stored))
			CreateFile(stored)
			df.Link()

			Expect(df.IsLinked()).To(BeTrue())
			Expect(df.Restore()).To(Succeed())

			Expect(Exists("config")).To(BeFalse())
			Expect(df.IsLinked()).To(BeFalse())
			Expect(df.IsReadyToBeStored()).To(BeTrue())
		})

		It("fails if file is not stored and linked properly", func() {
			CreateFile(stored())
			Expect(df().Restore()).NotTo(Succeed())
		})
	})

	Describe("Delete", func() {
		It("removes both stored file and link to it from original dotfile location", func() {
			o, s := orig(), stored()

			CreateFile(s)
			df().Link()

			Expect(df().Delete()).To(Succeed())

			Expect(Exists(s)).To(BeFalse())
			Expect(Exists(o)).To(BeFalse())
		})

		It("removes empty nested directories in both places", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			CreateDir(filepath.Dir(stored))
			CreateFile(stored)
			df.Link()

			Expect(df.IsLinked()).To(BeTrue())

			df.Delete()

			Expect(Exists("config")).To(BeFalse())
			Expect(Exists(".config")).To(BeFalse())
		})

		It("fails if file is not stored and linked properly", func() {
			Expect(df().Delete()).NotTo(Succeed())
		})
	})
})
