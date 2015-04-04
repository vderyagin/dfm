package dotfile_test

import (
	"log"
	"path/filepath"

	. "gitlab.com/vderyagin/dfm/dotfile"
	. "gitlab.com/vderyagin/dfm/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CurrentState", func() {
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

	It("correctly assigns 'Linked' state", func() {
		CreateFile(orig())
		df().Store()

		Expect(df().CurrentState().String()).To(Equal(Linked.String()))
	})

	It("correctly assigns 'NotLinked' state", func() {
		CreateFile(stored())

		Expect(df().CurrentState().String()).To(Equal(NotLinked.String()))
	})

	It("correctly assigns 'Conflict' state", func() {
		CreateFile(stored())
		CreateFile(orig())

		Expect(df().CurrentState().String()).To(Equal(Conflict.String()))
	})

	It("correctly assigns 'Missing' state", func() {
		Expect(df().CurrentState().String()).To(Equal(Missing.String()))
	})
})
