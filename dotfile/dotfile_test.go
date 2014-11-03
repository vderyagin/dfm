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
	Describe("IsStored", func() {
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

		It("Returns true if file is properly stored", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink(stored(), orig())

			Expect(df().IsStored()).To(BeTrue())
		})

		It("Returns false if original file location is empty", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)

			Expect(df().IsStored()).To(BeFalse())
		})

		It("Returns false if stored file location is empty", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink(stored(), orig())
			os.Remove(stored())

			Expect(df().IsStored()).To(BeFalse())
		})

		It("Returns false if original file is not a symlink", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			ioutil.WriteFile(orig(), []byte{}, 0777)

			Expect(df().IsStored()).To(BeFalse())
		})

		It("Returns false if stored file is not a regular file", func() {
			os.MkdirAll(stored(), 0777)
			os.Symlink(stored(), orig())

			Expect(df().IsStored()).To(BeFalse())
		})

		It("Returns false both files are not linked properly", func() {
			ioutil.WriteFile(stored(), []byte{}, 0777)
			os.Symlink("/wrong/location", orig())

			Expect(df().IsStored()).To(BeFalse())
		})
	})
})
