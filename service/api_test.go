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

var _ = Describe("testing api functionality", func() {
	Describe("testing function NewCache", newCacheTest)
	Describe("testing function hashKeyToInt", hashKeyToIntTest)
})

func newCacheTest() {
	Context("Given a valid K data type and a valid setSize", func() {
		It("should return a not nil instance with not empty fields, setSize = 2 (provided as input) and a nil error", func() {
			cache, err := NewCache[int, string](2)
			Expect(err).Error().ShouldNot(HaveOccurred())

			Expect(cache).ShouldNot(BeNil())
			Expect(cache.sets).ShouldNot(BeNil())
			Expect(cache.setSize).Should(Equal(2))
			Expect(cache.setSize).ShouldNot(BeNil())
		})
	})

	Context("Given a setSize = 0", func() {
		It("should return an error cause it can't be used for mod functionality", func() {
			Expect(NewCache[int, string](0)).Error().Should(HaveOccurred())
		})
	})

	Context("Given a setSize < 0", func() {
		It("should return an error cause it can't be used for mod functionality", func() {
			Expect(NewCache[int, string](-1)).Error().Should(HaveOccurred())
		})
	})

	Context("Busniess rule: Given K data type is not a primitive data type", func() {
		Context("Given an structure", func() {
			It("should return an error error", func() {
				Expect(NewCache[struct{}, string](4)).Error().Should(HaveOccurred())
			})
		})

		Context("Given a pointer data type", func() {
			It("should return an error error", func() {
				Expect(NewCache[*struct{}, string](4)).Error().Should(HaveOccurred())
			})
		})
	})
}

func hashKeyToIntTest() {
	When("Testing int as input values", func() {
		functionInvoker := new(hashKeyToIntImpl[int])
		Context("Given 1 (int) as input", func() {
			It("should return '1645622373'", func() {
				Expect(functionInvoker.hashKeyToInt(1)).Should(Equal(1645622373))
			})
		})

		Context("Given 2 (int) as input", func() {
			It("should return '517158638'", func() {
				Expect(functionInvoker.hashKeyToInt(2)).Should(Equal(517158638))
			})
		})

		Context("Given a maxInt value as input", func() {
			It("should return '2756140996'", func() {
				Expect(functionInvoker.hashKeyToInt(math.MaxInt)).Should(Equal(2756140996))
			})
		})
	})

	When("Testing case sensitive cases for string data type", func() {
		functionInvoker := new(hashKeyToIntImpl[string])

		Context("Given 'foo' (string) as input", func() {
			It("should return '2038559810'", func() {
				Expect(functionInvoker.hashKeyToInt("foo")).Should(Equal(2038559810))
			})
		})

		Context("Given 'FOO' (string) as input", func() {
			It("should return '3571895650'", func() {
				Expect(functionInvoker.hashKeyToInt("FOO")).Should(Equal(3571895650))
			})
		})

		Context("Given 'Foo' as input and 'foo'", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := functionInvoker.hashKeyToInt("Foo") == functionInvoker.hashKeyToInt("foo")
				Expect(comparissionResult).Should(BeFalse())
			})
		})
	})

	When("Testing same input value, but with different data type", func() {
		functionInvoker := new(hashKeyToIntImpl[any])

		Context("Given '1' (int) as input and '1' (string)", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := functionInvoker.hashKeyToInt("1") == functionInvoker.hashKeyToInt(1)
				Expect(comparissionResult).Should(BeFalse())
			})
		})

		Context("Given 'true' (string) as input and 'true' (bool)", func() {
			It("should return FALSE if compare their results", func() {
				comparissionResult := functionInvoker.hashKeyToInt("true") == functionInvoker.hashKeyToInt(true)
				Expect(comparissionResult).Should(BeFalse())
			})
		})
	})

	When("Testing same input value and same data type", func() {
		functionInvoker := new(hashKeyToIntImpl[any])

		Context("Given the same input for string data type in different calls", func() {
			It("should return TRUE if compare their results", func() {
				comparissionResult := functionInvoker.hashKeyToInt("foo") == functionInvoker.hashKeyToInt("foo")
				Expect(comparissionResult).Should(BeTrue())
			})
		})

		Context("Given the same input for int data type in different calls", func() {
			It("should return TRUE if compare their results", func() {
				comparissionResult := functionInvoker.hashKeyToInt(123) == functionInvoker.hashKeyToInt(123)
				Expect(comparissionResult).Should(BeTrue())
			})
		})
	})
}
