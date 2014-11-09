package fsutil_test

import (
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/fsutil"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FSutil", func() {
	ExecuteEachInTempDir()

	Describe("FilesIn", func() {
		It("returns empty collection if argument does not exist", func() {
			Expect(FilesIn("nonexistent_dir")).To(BeEmpty())
		})

		It("returns empty collection for empty directory", func() {
			Expect(FilesIn(".")).To(BeEmpty())
		})

		It("returns empty collection for directory with hidden files", func() {
			CreateFile(".hidden")
			CreateFile(".foo/bar/baz")

			Expect(FilesIn(".")).To(BeEmpty())
		})

		It("returns collection including regular non-hidden files", func() {
			baseName := "bashrc"
			absPath, _ := filepath.Abs(baseName)
			CreateFile(baseName)

			Expect(FilesIn(".")).To(ContainElement(absPath))
		})

		It("returns collection including deeply nested files", func() {
			relPath := "config/camlistore/server-config.json"
			absPath, _ := filepath.Abs(relPath)
			CreateDir(filepath.Dir(absPath))
			CreateFile(relPath)

			Expect(FilesIn(".")).To(ContainElement(absPath))
		})
	})

	Describe("IsEmptyDir", func() {
		It("returns true if path is an empty directory", func() {
			Expect(IsEmptyDir(".")).To(BeTrue())
		})

		It("returns false if path contains stuff", func() {
			CreateDir("some_dir")
			Expect(IsEmptyDir(".")).To(BeFalse())
		})

		It("returns false if path does not exist", func() {
			Expect(IsEmptyDir("/non/existent/directory")).To(BeFalse())
		})

		It("returns false if path is not a directory", func() {
			CreateFile("file")
			Expect(IsEmptyDir("file")).To(BeFalse())
		})
	})

	Describe("DeleteEmptyDirs", func() {
		It("deletes all directories in hierarchy until non-empty one", func() {
			CreateDir("foo/bar/baz/quux")
			CreateFile("foo/a_file")

			Expect(DeleteEmptyDirs("foo/bar/baz/quux")).To(Succeed())
			Expect(Exists("foo/bar")).To(BeFalse())
			Expect(Exists("foo")).To(BeTrue())
		})

		It("silengly exits if directory does not exist", func() {
			Expect(DeleteEmptyDirs("non_existent_directory")).To(Succeed())
		})

		It("silengly exits if argument is not a directory", func() {
			CreateFile("a_file")
			Expect(DeleteEmptyDirs("a_file")).To(Succeed())
			Expect(Exists("a_file")).To(BeTrue())
		})
	})

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
})
