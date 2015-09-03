package pezinventory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET /v1/inventory", func() {
	Context("When true", func() {
		It("should return true", func() {
			Ω(true).To(BeTrue())
		})
	})
})
