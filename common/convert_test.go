package common

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"unsafe"
)

var _ = Describe("Convert", func() {
	Context("StrToBytes", func() {
		It("normal", func() {
			s := "hello world"
			b := StrToBytes(s)

			bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
			sHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
			Expect(bHeader.Data == sHeader.Data).Should(BeTrue())
			Expect(string(b) == "hello world").Should(BeTrue())
		})

		It("nil", func() {
			var s string
			b := StrToBytes(s)

			bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
			sHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
			Expect(bHeader.Data == sHeader.Data).Should(BeTrue())
			Expect(string(b) == "").Should(BeTrue())
		})
	})

	Context("BytesToStr", func() {
		It("normal", func() {
			b := []byte("hello world")
			s := BytesToStr(b)

			bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
			sHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
			Expect(bHeader.Data == sHeader.Data).Should(BeTrue())
			Expect(bytes.Compare(b,[]byte("hello world")) == 0).Should(BeTrue())
		})

		It("nil", func() {
			var b []byte
			s := BytesToStr(b)
			
			bHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
			sHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))
			Expect(bHeader.Data == sHeader.Data).Should(BeTrue())
			Expect(bytes.Compare(b,[]byte("hello world")) == 0).Should(BeTrue())
		})
	})
})
