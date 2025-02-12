package cache

import (
	"container/list"
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
	Describe("testing function Put", putTest)
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

func putTest() {
	When("Testing 4-way with a single set scenario", func() {
		var (
			mockedHashKeyToIntConverter *hashKeyToIntConverterMock[int]
			cacheTester                 *Cache[int, any]
		)

		BeforeEach(func() {
			mockedHashKeyToIntConverter = new(hashKeyToIntConverterMock[int])
			cacheTester = &Cache[int, any]{
				setSize:               4,
				sets:                  make(map[int]*list.List),
				entries:               make(map[int]*list.Element),
				hashKeyToIntConverter: mockedHashKeyToIntConverter,
			}
		})

		Context("Given 4 invocations with the same input", func() {
			It("should have a single set defined in cache structure with a single element", func() {
				mockedHashKeyToIntConverter.On("hashKeyToInt", 123).Return(1)

				cacheTester.Put(123, "fooo")
				cacheTester.Put(123, "fooo")
				cacheTester.Put(123, "fooo")
				cacheTester.Put(123, "fooo")

				Expect(cacheTester.sets[1].Len()).Should(Equal(1))

				// Call Front or Back is the same cause it has just a single element
				cachedItem := cacheTester.sets[1].Front().Value.(*entry[int, any])
				Expect(cachedItem.key).Should(Equal(123))
				Expect(cachedItem.value).Should(Equal("fooo"))
			})
		})

		Context("Given 4 invocations with the same key and different value", func() {
			It("should have a single defined set, and the value stored in it should be the last saved value", func() {
				mockedHashKeyToIntConverter.On("hashKeyToInt", 123).Return(1)

				cacheTester.Put(123, "firstValue")
				cacheTester.Put(123, "secondValue")
				cacheTester.Put(123, "thirdValue")
				cacheTester.Put(123, "fourthValue")

				Expect(cacheTester.sets[1].Len()).Should(Equal(1))

				// Call Front or Back is the same cause it has just a single element
				cachedItem := cacheTester.sets[1].Front().Value.(*entry[int, any])
				Expect(cachedItem.key).Should(Equal(123))
				Expect(cachedItem.value).Should(Equal("fourthValue"))
			})
		})

		Context("Given 6 invocations with different key-value pairs", func() {
			Context("Given the same setIndex for all invocations", func() {
				It("should return have a single defined set,the last stored items", func() {
					mockedHashKeyToIntConverter.On("hashKeyToInt", 123).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 456).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 789).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 100).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 101).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 102).Return(0)

					cacheTester.Put(123, "firstValue")
					cacheTester.Put(456, "secondValue")
					cacheTester.Put(789, "thirdValue")
					cacheTester.Put(100, "fourthValue")
					cacheTester.Put(101, "fifthValue")
					cacheTester.Put(102, "sixthValue")

					Expect(cacheTester.sets[0].Len()).Should(Equal(4))

					// creates a list of validations
					expectedTestValues := []struct {
						setIndex          int
						expectedSetLength int
						expectedKey       int
						expectedValue     any
					}{
						{
							setIndex:          0,
							expectedSetLength: 4,
							expectedKey:       102,
							expectedValue:     "sixthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 3,
							expectedKey:       101,
							expectedValue:     "fifthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       100,
							expectedValue:     "fourthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       789,
							expectedValue:     "thirdValue",
						},
					}

					for _, expectedItem := range expectedTestValues {
						cachedSetItems := cacheTester.sets[expectedItem.setIndex]
						Expect(cachedSetItems.Len()).Should(Equal(expectedItem.expectedSetLength))

						// The item to be validated is being removed from the set in order to validate the next item configured as expected
						cachedItem := cachedSetItems.Remove(cachedSetItems.Front()).(*entry[int, any])
						Expect(cachedItem.key).Should(Equal(expectedItem.expectedKey))
						Expect(cachedItem.value).Should(Equal(expectedItem.expectedValue))
					}
				})
			})
		})

		Context("Given 4 invocations with different key-value pairs", func() {
			Context("Given a different set per invocation", func() {
				It("should have a 4 defined sets, with a single value per set", func() {
					mockedHashKeyToIntConverter.On("hashKeyToInt", 1).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 2).Return(1)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 3).Return(2)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 4).Return(3)

					cacheTester.Put(1, true)
					cacheTester.Put(2, "hello")
					cacheTester.Put(3, 123)
					cacheTester.Put(4, false)

					// creates a list of validations
					expectedTestValues := []struct {
						setIndex          int
						expectedSetLength int
						expectedKey       int
						expectedValue     any
					}{
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       1,
							expectedValue:     true,
						},
						{
							setIndex:          1,
							expectedSetLength: 1,
							expectedKey:       2,
							expectedValue:     "hello",
						},
						{
							setIndex:          2,
							expectedSetLength: 1,
							expectedKey:       3,
							expectedValue:     123,
						},
						{
							setIndex:          3,
							expectedSetLength: 1,
							expectedKey:       4,
							expectedValue:     false,
						},
					}

					for _, expectedItem := range expectedTestValues {
						cachedSetItems := cacheTester.sets[expectedItem.setIndex]
						Expect(cachedSetItems.Len()).Should(Equal(expectedItem.expectedSetLength))

						cachedItem := cachedSetItems.Front().Value.(*entry[int, any])
						Expect(cachedItem.key).Should(Equal(expectedItem.expectedKey))
						Expect(cachedItem.value).Should(Equal(expectedItem.expectedValue))
					}
				})
			})

			Context("Given the same set twice", func() {
				It("should have 2 defined sets, with a two values per set", func() {
					mockedHashKeyToIntConverter.On("hashKeyToInt", 1).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 2).Return(1)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 3).Return(0)
					mockedHashKeyToIntConverter.On("hashKeyToInt", 4).Return(1)

					cacheTester.Put(1, true)
					cacheTester.Put(2, "hello")
					cacheTester.Put(3, 123)
					cacheTester.Put(4, false)

					// creates a list of validations
					expectedTestValues := []struct {
						setIndex          int
						expectedSetLength int
						expectedKey       int
						expectedValue     any
					}{
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       3,
							expectedValue:     123,
						},
						{
							setIndex:          1,
							expectedSetLength: 2,
							expectedKey:       4,
							expectedValue:     false,
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       1,
							expectedValue:     true,
						},
						{
							setIndex:          1,
							expectedSetLength: 1,
							expectedKey:       2,
							expectedValue:     "hello",
						},
					}

					for _, expectedItem := range expectedTestValues {
						cachedSetItems := cacheTester.sets[expectedItem.setIndex]
						Expect(cachedSetItems.Len()).Should(Equal(expectedItem.expectedSetLength))

						// The item to be validated is being removed from the set in order to validate the next item configured as expected
						cachedItem := cachedSetItems.Remove(cachedSetItems.Front()).(*entry[int, any])
						Expect(cachedItem.key).Should(Equal(expectedItem.expectedKey))
						Expect(cachedItem.value).Should(Equal(expectedItem.expectedValue))
					}
				})
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
