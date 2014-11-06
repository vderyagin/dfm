package fsutil_test

import (
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

			Expect(DeleteEmptyDirs("foo/bar/baz/quux")).To(BeNil())
			Expect(Exists("foo/bar")).To(BeFalse())
			Expect(Exists("foo")).To(BeTrue())
		})

		It("silengly exits if directory does not exist", func() {
			Expect(DeleteEmptyDirs("non_existent_directory")).To(BeNil())
		})

		It("silengly exits if argument is not a directory", func() {
			CreateFile("a_file")
			Expect(DeleteEmptyDirs("a_file")).To(BeNil())
			Expect(Exists("a_file")).To(BeTrue())
		})
	})
})
