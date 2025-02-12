package cache

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestMainSuite
func TestMainSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("testing hashKeyToInt functionality", hashKeyToIntTest)

func hashKeyToIntTest() {
	When("Testing int as input values", func() {
		Context("Given 1 (int) as input", func() {
			It("should return '873244444'", func() {
				Expect(hashKeyToInt(1)).Should(Equal(873244444))
			})
		})

		Context("Given 2 (int) as input", func() {
			It("should return '923577301'", func() {
				Expect(hashKeyToInt(2)).Should(Equal(923577301))
			})
		})

		Context("Given a maxInt value as input", func() {
			It("should return '546809323'", func() {
				Expect(hashKeyToInt(math.MaxInt)).Should(Equal(546809323))
			})
		})
	})

	When("Testing case sensitive cases for string data type", func() {
		Context("Given 'foo' (string) as input", func() {
			It("should return '2851307223'", func() {
				Expect(hashKeyToInt("foo")).Should(Equal(2851307223))
			})
		})

		Context("Given 'FOO' (string) as input", func() {
			It("should return '3972896247'", func() {
				Expect(hashKeyToInt("FOO")).Should(Equal(3972896247))
			})
		})

		Context("Given 'Foo' as input and 'foo'", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := hashKeyToInt("Foo") == hashKeyToInt("foo")
				Expect(comparissionResult).Should(BeFalse())
			})
		})
	})

	// TODO: Implement functionality that allows to meet this test scenario
	When("Testing same input value, but with different data type", func() {
		Context("Given '1' (int) as input and '1' (string)", func() {
			It("should return FALSE if compare their results", func() {
				result := hashKeyToInt("1") == hashKeyToInt(1)
				Expect(result).Should(BeFalse())
			})
		})
	})
}
