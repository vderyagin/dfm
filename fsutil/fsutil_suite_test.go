package fsutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFsutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fsutil Suite")
}
