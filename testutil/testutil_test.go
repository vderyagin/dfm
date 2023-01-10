package testutil_test

import (
	. "github.com/vderyagin/dfm/fsutil"
	. "github.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestUtil", func() {
	ExecuteEachInTempDir()

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
