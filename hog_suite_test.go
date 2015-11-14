package hog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShahog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hog Suite")
}
