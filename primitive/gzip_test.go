package primitive

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gzip", func() {
	Context("gzip & gunzip", func() {
		It("should be succeed", func() {
			update := []byte{2, 5, 234, 163, 145, 177, 10, 0, 136, 247, 203, 185, 212, 9, 1, 2, 119, 1, 99, 119, 1, 100, 39, 1, 5, 117, 115, 101, 114, 115, 36, 51, 50, 53, 97, 97, 48, 52, 54, 45, 100, 99, 97, 101, 45, 52, 54, 102, 48, 45, 97, 56, 57, 56, 45, 49, 99, 102, 51, 102, 49, 49, 102, 98, 54, 51, 49, 1, 39, 0, 234, 163, 145, 177, 10, 2, 3, 105, 100, 115, 0, 39, 0, 234, 163, 145, 177, 10, 2, 2, 100, 115, 0, 8, 0, 234, 163, 145, 177, 10, 3, 1, 123, 65, 228, 196, 138, 61, 64, 0, 0, 5, 247, 203, 185, 212, 9, 0, 8, 1, 6, 98, 108, 111, 99, 107, 115, 2, 119, 1, 97, 119, 1, 98, 39, 1, 5, 117, 115, 101, 114, 115, 36, 49, 50, 56, 98, 99, 56, 48, 49, 45, 52, 48, 98, 52, 45, 52, 53, 54, 97, 45, 98, 100, 101, 51, 45, 50, 56, 51, 99, 57, 57, 97, 52, 98, 101, 49, 97, 1, 39, 0, 247, 203, 185, 212, 9, 2, 3, 105, 100, 115, 0, 39, 0, 247, 203, 185, 212, 9, 2, 2, 100, 115, 0, 8, 0, 247, 203, 185, 212, 9, 3, 1, 123, 65, 227, 81, 204, 190, 224, 0, 0, 0}
			data, err := Gzip(update)
			Expect(err).Should(Succeed())

			update2, err := Gunzip(data)
			Expect(err).Should(Succeed())

			Expect(bytes.Compare(update, update2) == 0).Should(BeTrue())
		})

		It("gzip nil", func() {
			data, err := Gzip(nil)
			Expect(err).Should(Succeed())

			data2, err := Gunzip(data)
			Expect(err).Should(Succeed())

			Expect(bytes.Compare(data, data2) == 0).Should(BeTrue())
		})
	})
})
