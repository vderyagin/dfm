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

	Describe("SymlinksIn", func() {
		It("returns empty closed channel if argument does not exist", func() {
			Expect(chanToSlice(SymlinksIn("nonexistent_dir"))).To(BeEmpty())
		})

		It("returns empty closed channel for empty directory", func() {
			Expect(chanToSlice(SymlinksIn("."))).To(BeEmpty())
		})

		It("returns empty closed channel for directory with only regular files", func() {
			CreateFile("foo")
			CreateFile("bar/baz")

			Expect(chanToSlice(SymlinksIn("."))).To(BeEmpty())
		})

		It("returns empty closed channel for directory with hidden symlinks", func() {
			CreateFile("target")
			os.Symlink("target", ".hidden_link")

			Expect(chanToSlice(SymlinksIn("."))).To(BeEmpty())
		})

		It("returns channel producing symlinks", func() {
			CreateFile("target")
			os.Symlink("target", "link")
			absPath, _ := filepath.Abs("link")

			Expect(chanToSlice(SymlinksIn("."))).To(ContainElement(absPath))
		})

		It("returns channel producing nested symlinks", func() {
			CreateFile("target")
			CreateDir("nested/dir")
			os.Symlink("../../target", "nested/dir/link")
			absPath, _ := filepath.Abs("nested/dir/link")

			Expect(chanToSlice(SymlinksIn("."))).To(ContainElement(absPath))
		})

		It("does not return regular files", func() {
			CreateFile("regular")
			CreateFile("target")
			os.Symlink("target", "link")

			results := chanToSlice(SymlinksIn("."))
			Expect(results).To(HaveLen(1))
		})
	})

	Describe("IsRelativeSymlinkWithinDir", func() {
		It("returns true for relative symlink pointing within dir", func() {
			CreateFile("target")
			os.Symlink("target", "link")
			dir, _ := filepath.Abs(".")

			Expect(IsRelativeSymlinkWithinDir("link", dir)).To(BeTrue())
		})

		It("returns true for relative symlink in nested directory", func() {
			CreateFile("target")
			CreateDir("nested")
			os.Symlink("../target", "nested/link")
			dir, _ := filepath.Abs(".")

			Expect(IsRelativeSymlinkWithinDir("nested/link", dir)).To(BeTrue())
		})

		It("returns false for absolute symlink", func() {
			CreateFile("target")
			absTarget, _ := filepath.Abs("target")
			os.Symlink(absTarget, "link")
			dir, _ := filepath.Abs(".")

			Expect(IsRelativeSymlinkWithinDir("link", dir)).To(BeFalse())
		})

		It("returns false for symlink pointing outside dir", func() {
			CreateDir("inside")
			CreateFile("outside")
			os.Symlink("../outside", "inside/link")
			dir, _ := filepath.Abs("inside")

			Expect(IsRelativeSymlinkWithinDir("inside/link", dir)).To(BeFalse())
		})

		It("returns false for non-symlink file", func() {
			CreateFile("regular")
			dir, _ := filepath.Abs(".")

			Expect(IsRelativeSymlinkWithinDir("regular", dir)).To(BeFalse())
		})

		It("returns false for non-existent path", func() {
			dir, _ := filepath.Abs(".")

			Expect(IsRelativeSymlinkWithinDir("nonexistent", dir)).To(BeFalse())
		})
	})

	Describe("ResolveSymlink", func() {
		It("returns absolute path for relative symlink", func() {
			CreateFile("target")
			os.Symlink("target", "link")
			expected, _ := filepath.Abs("target")

			resolved, err := ResolveSymlink("link")
			Expect(err).NotTo(HaveOccurred())
			Expect(resolved).To(Equal(expected))
		})

		It("returns absolute path for nested relative symlink", func() {
			CreateFile("target")
			CreateDir("nested")
			os.Symlink("../target", "nested/link")
			expected, _ := filepath.Abs("target")

			resolved, err := ResolveSymlink("nested/link")
			Expect(err).NotTo(HaveOccurred())
			Expect(resolved).To(Equal(expected))
		})

		It("returns cleaned path for absolute symlink", func() {
			CreateFile("target")
			absTarget, _ := filepath.Abs("target")
			os.Symlink(absTarget, "link")

			resolved, err := ResolveSymlink("link")
			Expect(err).NotTo(HaveOccurred())
			Expect(resolved).To(Equal(absTarget))
		})

		It("returns error for non-symlink", func() {
			CreateFile("regular")

			_, err := ResolveSymlink("regular")
			Expect(err).To(HaveOccurred())
		})

		It("returns error for non-existent path", func() {
			_, err := ResolveSymlink("nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})
})
