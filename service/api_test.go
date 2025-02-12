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
			It("should return '1645622373'", func() {
				Expect(hashKeyToInt(1)).Should(Equal(1645622373))
			})
		})

		Context("Given 2 (int) as input", func() {
			It("should return '517158638'", func() {
				Expect(hashKeyToInt(2)).Should(Equal(517158638))
			})
		})

		Context("Given a maxInt value as input", func() {
			It("should return '2756140996'", func() {
				Expect(hashKeyToInt(math.MaxInt)).Should(Equal(2756140996))
			})
		})
	})

	When("Testing case sensitive cases for string data type", func() {
		Context("Given 'foo' (string) as input", func() {
			It("should return '2038559810'", func() {
				Expect(hashKeyToInt("foo")).Should(Equal(2038559810))
			})
		})

		Context("Given 'FOO' (string) as input", func() {
			It("should return '3571895650'", func() {
				Expect(hashKeyToInt("FOO")).Should(Equal(3571895650))
			})
		})

		Context("Given 'Foo' as input and 'foo'", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := hashKeyToInt("Foo") == hashKeyToInt("foo")
				Expect(comparissionResult).Should(BeFalse())
			})
		})
	})

	When("Testing same input value, but with different data type", func() {
		Context("Given '1' (int) as input and '1' (string)", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := hashKeyToInt("1") == hashKeyToInt(1)
				Expect(comparissionResult).Should(BeFalse())
			})
		})

		Context("Given 'true' (string) as input and 'true' (bool)", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := hashKeyToInt("true") == hashKeyToInt(true)
				Expect(comparissionResult).Should(BeFalse())
			})
		})
	})

	When("Testing same input value and same data type", func() {
		Context("Given the same input for string data type in different calls", func() {
			It("should return TRUE if compare their results", func() {
				comparissionResult := hashKeyToInt("foo") == hashKeyToInt("foo")
				Expect(comparissionResult).Should(BeTrue())
			})
		})

		Context("Given the same input for int data type in different calls", func() {
			It("should return TRUE if compare their results", func() {
				comparissionResult := hashKeyToInt(123) == hashKeyToInt(123)
				Expect(comparissionResult).Should(BeTrue())
			})
		})
	})
}
