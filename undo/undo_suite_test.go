package undo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUndo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Undo Suite")
}
