package pezinventory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/pezdispenser"
)

var _ = Describe("Random Controller", func() {
	Context("when called with some random arguments", func() {
		BeforeEach(func() {
			RandomController("hi there")
		})

		It("Should do something great", func() {
			Ω(true).Should(BeTrue())
		})
	})
})
