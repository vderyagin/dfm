package testutil_test

import (
	"os"

	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUtil", func() {
	ExecuteEachInTempDir()

	Describe("Exists", func() {
		It("returns true if given path corresponds to existing file", func() {
			CreateFile("a_file")
			Expect(Exists("a_file")).To(BeTrue())
		})

		It("returns true if given path corresponds to existing directory", func() {
			CreateDir("a_dir")
			Expect(Exists("a_dir")).To(BeTrue())
		})

		It("returns false if there's no file at given path", func() {
			Expect(Exists("nonexistent_file")).To(BeFalse())
		})

		It("returns true if there's symlink with missing target at given path", func() {
			CreateFile("a_file")
			os.Symlink("a_file", "a_symlink")
			os.Remove("a_file")
			Expect(Exists("a_symlink")).To(BeTrue())
		})
	})

	Describe("CreateDir", func() {
		dir := "foo/bar/baz"

		It("creates a directory with all needed parents", func() {
			CreateDir(dir)

			Expect(Exists(dir)).To(BeTrue())
		})

		It("exits silently when there's already directory at a given path", func() {
			CreateDir(dir)
			CreateDir(dir)
		})
	})

	Describe("CreateFile", func() {
		file := "foo/bar/baz"

		It("creates a directory with all needed parents", func() {
			CreateFile(file)

			Expect(Exists(file)).To(BeTrue())
		})

		It("exits silently when there's already a file at a given path", func() {
			CreateFile(file)
			CreateFile(file)
		})
	})
})
