package dotfile_test

import (
	"io/ioutil"
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
			ioutil.WriteFile(stored(), []byte{}, 0777)

			Expect(df().IsStored()).To(BeTrue())
		})

		It("returns false if stored file location is empty", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink(stored(), orig())
			os.Remove(stored())

			Expect(df().IsStored()).To(BeFalse())
		})

		It("returns false if stored file is not a regular file", func() {
			os.MkdirAll(stored(), 0777)
			os.Symlink(stored(), orig())

			Expect(df().IsStored()).To(BeFalse())
		})
	})

	Describe("IsLinked", func() {
		It("returns true if file is stored and linked properly", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink(stored(), orig())

			Expect(df().IsLinked()).To(BeTrue())
		})

		It("returns false if file is not stored properly", func() {
			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's nothing at original location", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)

			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's no symlink at original location", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.MkdirAll(orig(), 0777)

			Expect(df().IsLinked()).To(BeFalse())
		})

		It("returns false if there's wrong symlink at original location", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			ioutil.WriteFile("wrong_file", []byte{}, 0777)
			os.Symlink("wrong_file", orig())

			Expect(df().IsLinked()).To(BeFalse())
		})
	})

	Describe("IsReadyToBeStored", func() {
		It("returns true if regular file not conflicting with stored ones", func() {
			ioutil.WriteFile(orig(), []byte{}, 0777)

			Expect(df().IsReadyToBeStored()).To(BeTrue())
		})

		It("returns false if original location is empty", func() {
			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})

		It("returns false if original location is not a regular file", func() {
			os.Mkdir(orig(), 0777)

			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})

		It("returns false if conflicts with already stored file", func() {
			ioutil.WriteFile(orig(), []byte{}, 0777)
			ioutil.WriteFile(stored(), []byte{}, 0777)

			Expect(df().IsReadyToBeStored()).To(BeFalse())
		})
	})

	Describe("Store", func() {

		It("stores file", func() {
			ioutil.WriteFile(orig(), []byte{}, 0777)

			Expect(df().IsStored()).To(BeFalse())
			Expect(df().Store()).To(BeNil())
			Expect(df().IsStored()).To(BeTrue())
		})

		It("creates intermediate directories for nested file", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			os.MkdirAll(filepath.Dir(orig), 0777)
			ioutil.WriteFile(orig, []byte{}, 0777)

			Expect(df.IsStored()).To(BeFalse())
			Expect(df.Store()).To(BeNil())
			Expect(df.IsStored()).To(BeTrue())
		})

		It("fails if file is not ready to be stored", func() {
			Expect(df().Store()).NotTo(BeNil())
			Expect(df().IsStored()).To(BeFalse())
		})
	})

	Describe("Link", func() {
		It("symlinks stored file to it's original location", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)

			Expect(df().IsLinked()).To(BeFalse())
			Expect(df().Link()).To(BeNil())
			Expect(df().IsLinked()).To(BeTrue())
		})

		It("creates deeply nested directories if necessary", func() {
			stored, _ := filepath.Abs("config/camlistore/server-config.json")
			orig, _ := filepath.Abs(".config/camlistore/server-config.json")
			df := New(stored, orig)
			os.MkdirAll(filepath.Dir(stored), 0777)
			ioutil.WriteFile(stored, []byte{}, 0777)

			Expect(df.IsLinked()).To(BeFalse())
			Expect(df.Link()).To(BeNil())
			Expect(df.IsLinked()).To(BeTrue())
		})

		It("fails if file is not stored", func() {
			Expect(df().Link()).NotTo(BeNil())
		})

		It("fails if file is already linked", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink(stored(), orig())

			Expect(df().Link()).NotTo(BeNil())
		})

		It("fails if there's conficting file at original location", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			ioutil.WriteFile(orig(), []byte{}, 0777)

			Expect(df().Link()).NotTo(BeNil())
		})
	})
})
