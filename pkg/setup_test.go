package pkg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJenkinsClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg test")
}
