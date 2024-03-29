package fsutil_test

import (
	"os"
	"path/filepath"

	. "github.com/vderyagin/dfm/fsutil"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func chanToSlice(input <-chan string) []string {
	output := []string{}

	for s := range input {
		output = append(output, s)
	}

	return output
}

var _ = Describe("FSutil", func() {
	ExecuteEachInTempDir()
	Describe("FilesIn", func() {
		It("returns empty closed channel if argument does not exist", func() {
			Expect(chanToSlice(FilesIn("nonexistent_dir"))).To(BeEmpty())
		})

		It("returns empty closed channel for empty directory", func() {
			Expect(chanToSlice(FilesIn("."))).To(BeEmpty())
		})

		It("returns empty closed channel for directory with hidden files", func() {
			CreateFile(".hidden")
			CreateFile(".foo/bar/baz")

			Expect(chanToSlice(FilesIn("."))).To(BeEmpty())
		})

		It("returns channel producing regular non-hidden files", func() {
			baseName := "bashrc"
			absPath, _ := filepath.Abs(baseName)
			CreateFile(baseName)

			Expect(chanToSlice(FilesIn("."))).To(ContainElement(absPath))
		})

		It("returns channel producing deeply nested files", func() {
			relPath := "config/camlistore/server-config.json"
			absPath, _ := filepath.Abs(relPath)
			CreateDir(filepath.Dir(absPath))
			CreateFile(relPath)

			Expect(chanToSlice(FilesIn("."))).To(ContainElement(absPath))
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

	Describe("IsRegularFile", func() {
		It("returns true for regular file", func() {
			CreateFile("foo/bar")
			Expect(IsRegularFile("foo/bar")).To(BeTrue())
		})

		It("returns false when file does not exist", func() {
			Expect(IsRegularFile("foo")).To(BeFalse())
		})

		It("returns false for directories", func() {
			CreateDir("foo")
			Expect(IsRegularFile("foo")).To(BeFalse())
		})
	})

	Describe("IsSymlink", func() {
		It("returns true for symlink", func() {
			CreateFile("foo")
			os.Symlink("foo", "bar")

			Expect(IsSymlink("bar")).To(BeTrue())
		})

		It("returns true for dangling symlink", func() {
			CreateFile("foo")
			os.Symlink("foo", "bar")
			os.Remove("foo")

			Expect(IsSymlink("bar")).To(BeTrue())
		})

		It("returns false for plain files", func() {
			CreateFile("foo")

			Expect(IsSymlink("foo")).To(BeFalse())
		})

		It("returns false for directories", func() {
			CreateDir("foo")

			Expect(IsSymlink("foo")).To(BeFalse())
		})
	})
})
