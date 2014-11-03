package dotfile_test

import (
	"io/ioutil"
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

		It("Returns true if file is properly stored", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")
			ioutil.WriteFile(stored, []byte{}, 0777)
			os.Symlink(stored, orig)

			Expect(New(stored, orig).IsStored()).To(BeTrue())
		})

		It("Returns false if original file location is empty", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")
			ioutil.WriteFile(stored, []byte{}, 0777)

			Expect(New(stored, orig).IsStored()).To(BeFalse())
		})

		It("Returns false if stored file location is empty", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")

			Expect(New(stored, orig).IsStored()).To(BeFalse())
		})

		It("Returns false if original file is not a symlink", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")
			ioutil.WriteFile(stored, []byte{}, 0777)
			ioutil.WriteFile(orig, []byte{}, 0777)

			Expect(New(stored, orig).IsStored()).To(BeFalse())
		})

		It("Returns false if stored file is not a regular file", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")
			os.MkdirAll(stored, 0777)

			Expect(New(stored, orig).IsStored()).To(BeFalse())
		})

		It("Returns false both files are not linked properly", func() {
			stored, _ := filepath.Abs("foo")
			orig, _ := filepath.Abs(".foo")
			ioutil.WriteFile(stored, []byte{}, 0777)
			os.Symlink("/wrong/location", orig)

			Expect(New(stored, orig).IsStored()).To(BeFalse())
		})
	})
})
