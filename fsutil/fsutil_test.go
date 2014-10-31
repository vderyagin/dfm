package fsutil_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/fsutil"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fsutil", func() {
	Describe("FilesIn", func() {
		ExecuteEachInTempDir()
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
