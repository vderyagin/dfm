package dotfile_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDotfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dotfile Suite")
}
