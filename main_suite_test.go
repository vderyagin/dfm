package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDfm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DFM Suite")
}
