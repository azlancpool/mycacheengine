package cache

import (
	"container/list"
	"math"
	"sync"
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
	Describe("testing function Get", getTest)
	Describe("testing function ListAll", listAllTest)
	Describe("testing function Delete", deleteTest)
	Describe("testing function hashKeyToInt", hashKeyToIntTest)
	Describe("testing concurrency", concurrencyTest)
})

type itemsToLoad struct {
	key   int
	value any
}

type itemsToLoadWithMockedKeyHashed struct {
	key             int
	value           any
	mockedHashedKey int
}

func newCacheTest() {
	Context("Given a valid K data type and a valid setSize", func() {
		It(`should return:
			- A not nil instance
			- Instance fields must not be empty
			- setSize = 2 (provided as input)
			- Not nil function getItemToRemove (MRU scenario)
			- Nil error`, func() {
			cache, err := NewCache[int, string](2, MRU_ALGO)
			Expect(err).Error().ShouldNot(HaveOccurred())

			Expect(cache).ShouldNot(BeNil())
			Expect(cache.sets).ShouldNot(BeNil())
			Expect(cache.setSize).Should(Equal(2))
			Expect(cache.setSize).ShouldNot(BeNil())
			Expect(cache.getItemToRemove).ShouldNot(BeNil())
		})

		It(`should return:
			- A not nil instance
			- Instance fields must not be empty
			- setSize = 4 (provided as input)
			- Not nil function getItemToRemove (LRU scenario)
			- Nil error`, func() {
			cache, err := NewCache[int, string](4, LRU_ALGO)
			Expect(err).Error().ShouldNot(HaveOccurred())

			Expect(cache).ShouldNot(BeNil())
			Expect(cache.sets).ShouldNot(BeNil())
			Expect(cache.setSize).Should(Equal(4))
			Expect(cache.setSize).ShouldNot(BeNil())
			Expect(cache.getItemToRemove).ShouldNot(BeNil())
		})
	})

	Context("Given a nil replacement algorithm", func() {
		It(`should return:
			- A not nil instance
			- Instance fields must not be empty
			- setSize = 2 (provided as input)
			- Not nil function getItemToRemove (LRU scenario)
			- Nil error`, func() {
			cache, err := NewCache[int, string](2)
			Expect(err).Error().ShouldNot(HaveOccurred())

			Expect(cache).ShouldNot(BeNil())
			Expect(cache.sets).ShouldNot(BeNil())
			Expect(cache.setSize).Should(Equal(2))
			Expect(cache.setSize).ShouldNot(BeNil())
			Expect(cache.getItemToRemove).ShouldNot(BeNil())
		})
	})

	Context("Given a setSize = 0", func() {
		It("should return an error cause it can't be used for mod functionality", func() {
			Expect(NewCache[int, string](0, MRU_ALGO)).Error().Should(HaveOccurred())
		})
	})

	Context("Given a setSize < 0", func() {
		It("should return an error cause it can't be used for mod functionality", func() {
			Expect(NewCache[int, string](-1, MRU_ALGO)).Error().Should(HaveOccurred())
		})
	})

	Context("Busniess rule: Given K data type is not a primitive data type", func() {
		Context("Given an structure", func() {
			It("should return an error error", func() {
				Expect(NewCache[struct{}, string](4, MRU_ALGO)).Error().Should(HaveOccurred())
			})
		})

		Context("Given a pointer data type", func() {
			It("should return an error error", func() {
				Expect(NewCache[*struct{}, string](4, MRU_ALGO)).Error().Should(HaveOccurred())
			})
		})
	})
}

func putTest() {
	When("LRU - Testing 4-way-set-associative-cache", func() {
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
				getItemToRemove:       LRU_ITEM_TO_REMOVE_GETTER,
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
					itemsToLoad := []itemsToLoad{
						{123, "firstValue"},
						{456, "secondValue"},
						{789, "thirdValue"},
						{100, "fourthValue"},
						{101, "fifthValue"},
						{102, "sixthValue"},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
						cacheTester.Put(item.key, item.value)
					}

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

		Context("Given 8 invocations with different key-value pairs", func() {
			Context("Given the same setIndex for all invocations", func() {
				It("should return have a single defined set,the last stored items", func() {
					itemsToLoad := []itemsToLoad{
						{123, "firstValue"},
						{456, "secondValue"},
						{789, "thirdValue"},
						{100, "fourthValue"},
						{101, "fifthValue"},
						{102, "sixthValue"},
						{456, "anotherSecondValue"},
						{100, "anotherFourthValue"},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
						cacheTester.Put(item.key, item.value)
					}

					// Test scenario explanation
					// For the keys: [123, 456, 789, 100, 101, 102, 456, 100]
					// It should follow the iterations:
					// 123, 456, 789, 100 // First 4 iterations, set is full since here
					// 456, 789, 100, 101 // 101 is a new key, 123 is removed cause algo is LRU
					// 789, 100, 101, 102 // 102 is a new key, 456 is removed cause algo is LRU
					// 100, 101, 102, 456 // 456 is a new key, 456 is removed cause algo is LRU
					// 101, 102, 456, 100 // 100 is NOT a new key, 456 it's moved to the front of the set

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
							expectedKey:       100,
							expectedValue:     "anotherFourthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 3,
							expectedKey:       456,
							expectedValue:     "anotherSecondValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       102,
							expectedValue:     "sixthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       101,
							expectedValue:     "fifthValue",
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
					itemsToLoad := []itemsToLoadWithMockedKeyHashed{
						{1, true, 0},
						{2, "hello", 1},
						{3, 123, 2},
						{4, false, 3},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashedKey)
						cacheTester.Put(item.key, item.value)
					}

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
		})

		Context("Testing LRU - SHOW CASE defined in README input=[1,7,9,15,9,7,45]", func() {
			Context("Given the same set per key-value pair", func() {
				It("should have a single set, with 4 key-value pairs", func() {
					itemsToLoad := []itemsToLoad{
						{1, true},
						{7, "hello"},
						{9, 123},
						{15, false},
						{9, 456},
						{7, "foo"},
						{45, "bar"},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
						cacheTester.Put(item.key, item.value)
					}

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
							expectedKey:       45,
							expectedValue:     "bar",
						},
						{
							setIndex:          0,
							expectedSetLength: 3,
							expectedKey:       7,
							expectedValue:     "foo",
						},
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       9,
							expectedValue:     456,
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       15,
							expectedValue:     false,
						},
					}

					for _, expectedItem := range expectedTestValues {
						cachedSetItems := cacheTester.sets[expectedItem.setIndex]
						Expect(cachedSetItems.Len()).Should(Equal(expectedItem.expectedSetLength))

						cachedItem := cachedSetItems.Remove(cachedSetItems.Front()).(*entry[int, any])
						Expect(cachedItem.key).Should(Equal(expectedItem.expectedKey))
						Expect(cachedItem.value).Should(Equal(expectedItem.expectedValue))
					}
				})
			})

			Context("Given the same set twice", func() {
				It("should have 2 defined sets, with a two values per set", func() {
					itemsToLoad := []itemsToLoadWithMockedKeyHashed{
						{1, true, 0},
						{2, "hello", 1},
						{3, 123, 0},
						{4, false, 1},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashedKey)
						cacheTester.Put(item.key, item.value)
					}

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

	When("MRU - Testing 4-way-set-associative-cache", func() {
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
				getItemToRemove:       MRU_ITEM_TO_REMOVE_GETTER,
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
					itemsToLoad := []itemsToLoad{
						{123, "firstValue"},
						{456, "secondValue"},
						{789, "thirdValue"},
						{100, "fourthValue"},
						{101, "fifthValue"},
						{102, "sixthValue"},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
						cacheTester.Put(item.key, item.value)
					}

					// Test scenario explanation
					// For the keys: [123, 456, 789, 100, 101, 102, 456, 100]
					// It should follow the iterations:
					// 123, 456, 789, 100 // First 4 iterations, set is full since here
					// 123, 456, 789, 101 // 101 is a new key, 100 is removed cause algo is MRU
					// 123, 456, 789, 102 // 102 is a new key, 101 is removed cause algo is MRU

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
							expectedKey:       789,
							expectedValue:     "thirdValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       456,
							expectedValue:     "secondValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       123,
							expectedValue:     "firstValue",
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

		Context("Given 8 invocations with different key-value pairs", func() {
			Context("Given the same setIndex for all invocations", func() {
				It("should return have a single defined set,the last stored items", func() {
					itemsToLoad := []itemsToLoad{
						{123, "firstValue"},
						{456, "secondValue"},
						{789, "thirdValue"},
						{100, "fourthValue"},
						{101, "fifthValue"},
						{102, "sixthValue"},
						{456, "anotherSecondValue"},
						{100, "anotherFourthValue"},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
						cacheTester.Put(item.key, item.value)
					}

					// Test scenario explanation
					// For the keys: [123, 456, 789, 100, 101, 102, 456, 100]
					// It should follow the iterations:
					// 123, 456, 789, 100 // First 4 iterations, set is full since here
					// 123, 456, 789, 101 // 101 is a new key, 100 is removed cause algo is MRU
					// 123, 456, 789, 102 // 102 is a new key, 101 is removed cause algo is MRU
					// 123, 789, 102, 456 // 456 is NOT a new key, it's moved in front of the set.
					// 123, 789, 102, 100 // 100 is a new key, 456 is removed cause algo is MRU

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
							expectedKey:       100,
							expectedValue:     "anotherFourthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 3,
							expectedKey:       102,
							expectedValue:     "sixthValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 2,
							expectedKey:       789,
							expectedValue:     "thirdValue",
						},
						{
							setIndex:          0,
							expectedSetLength: 1,
							expectedKey:       123,
							expectedValue:     "firstValue",
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
					itemsToLoad := []itemsToLoadWithMockedKeyHashed{
						{1, true, 0},
						{2, "hello", 1},
						{3, 123, 2},
						{4, false, 3},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashedKey)
						cacheTester.Put(item.key, item.value)
					}

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
					itemsToLoad := []itemsToLoadWithMockedKeyHashed{
						{1, true, 0},
						{2, "hello", 1},
						{3, 123, 0},
						{4, false, 1},
					}
					for _, item := range itemsToLoad {
						mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashedKey)
						cacheTester.Put(item.key, item.value)
					}

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

			Context("Testing MRU - SHOW CASE defined in README input=[1,7,9,15,9,7,45]", func() {
				Context("Given the same set per key-value pair", func() {
					It("should have a single set, with 4 key-value pairs", func() {
						itemsToLoad := []itemsToLoad{
							{1, true},
							{7, "hello"},
							{9, 123},
							{15, false},
							{9, 456},
							{7, "foo"},
							{45, "bar"},
						}
						for _, item := range itemsToLoad {
							mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(0)
							cacheTester.Put(item.key, item.value)
						}

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
								expectedKey:       45,
								expectedValue:     "bar",
							},
							{
								setIndex:          0,
								expectedSetLength: 3,
								expectedKey:       9,
								expectedValue:     456,
							},
							{
								setIndex:          0,
								expectedSetLength: 2,
								expectedKey:       15,
								expectedValue:     false,
							},
							{
								setIndex:          0,
								expectedSetLength: 1,
								expectedKey:       1,
								expectedValue:     true,
							},
						}

						for _, expectedItem := range expectedTestValues {
							cachedSetItems := cacheTester.sets[expectedItem.setIndex]
							Expect(cachedSetItems.Len()).Should(Equal(expectedItem.expectedSetLength))

							cachedItem := cachedSetItems.Remove(cachedSetItems.Front()).(*entry[int, any])
							Expect(cachedItem.key).Should(Equal(expectedItem.expectedKey))
							Expect(cachedItem.value).Should(Equal(expectedItem.expectedValue))
						}
					})
				})
			})
		})
	})
}

// getTest tests Get functionality
func getTest() {
	When("Testing GET function", func() {
		var (
			mockedHashKeyToIntConverter *hashKeyToIntConverterMock[int]
			preloadedCache              *Cache[int, any]
		)

		BeforeEach(func() {
			mockedHashKeyToIntConverter = new(hashKeyToIntConverterMock[int])
			preloadedCache = &Cache[int, any]{
				setSize:               4,
				sets:                  make(map[int]*list.List),
				entries:               make(map[int]*list.Element),
				hashKeyToIntConverter: mockedHashKeyToIntConverter,
				getItemToRemove:       LRU_ITEM_TO_REMOVE_GETTER,
			}

			preloadedItems := []struct {
				key              int
				value            any
				mockedHashResult int
			}{
				{
					key:              123,
					value:            "foo",
					mockedHashResult: 0,
				},
				{
					key:              456,
					value:            "bar",
					mockedHashResult: 1,
				},
				{
					key:              789,
					value:            struct{}{},
					mockedHashResult: 0,
				},
				{
					key:              100,
					value:            true,
					mockedHashResult: 0,
				},
			}

			for _, item := range preloadedItems {
				mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashResult)
				preloadedCache.Put(item.key, item.value)
			}
		})

		Context("Looking for a key that exists: key=123", func() {
			It("should return true and value='foo'", func() {
				value, exist := preloadedCache.Get(123)
				Expect(exist).Should(BeTrue())
				Expect(value).Should(Equal("foo"))
			})
		})

		Context("Looking for a key that exists: key=100", func() {
			It("should return true and value=true", func() {
				value, exist := preloadedCache.Get(100)
				Expect(exist).Should(BeTrue())
				Expect(value).Should(Equal(true))
			})
		})

		Context("Looking for a key that exists: key=789", func() {
			It("should return true and an empty struct as result", func() {
				value, exist := preloadedCache.Get(789)
				Expect(exist).Should(BeTrue())
				Expect(value).Should(BeEquivalentTo(struct{}{}))
			})
		})

		Context("Looking for a key that exists: key=456", func() {
			It("should return true and value='bar'", func() {
				value, exist := preloadedCache.Get(456)
				Expect(exist).Should(BeTrue())
				Expect(value).Should(Equal("bar"))
			})
		})

		Context("Looking for a key that doesn't exist: key=999", func() {
			It("should return false and a nil value", func() {
				mockedHashKeyToIntConverter.On("hashKeyToInt", 999).Return(789)

				value, exist := preloadedCache.Get(999)
				Expect(exist).Should(BeFalse())
				Expect(value).Should(BeNil())
			})
		})
	})
}

// listAllTest tests ListAll functionality
func listAllTest() {
	When("Testing GET function", func() {
		var (
			mockedHashKeyToIntConverter *hashKeyToIntConverterMock[int]
			preloadedCache              *Cache[int, any]
			preloadedItems              = []struct {
				key              int
				value            any
				mockedHashResult int
			}{
				{
					key:              123,
					value:            "foo",
					mockedHashResult: 0,
				},
				{
					key:              456,
					value:            "bar",
					mockedHashResult: 1,
				},
				{
					key:              789,
					value:            struct{}{},
					mockedHashResult: 0,
				},
				{
					key:              100,
					value:            true,
					mockedHashResult: 0,
				},
			}
		)

		BeforeEach(func() {
			mockedHashKeyToIntConverter = new(hashKeyToIntConverterMock[int])
			preloadedCache = &Cache[int, any]{
				setSize:               4,
				sets:                  make(map[int]*list.List),
				entries:               make(map[int]*list.Element),
				hashKeyToIntConverter: mockedHashKeyToIntConverter,
				getItemToRemove:       LRU_ITEM_TO_REMOVE_GETTER,
			}

			for _, item := range preloadedItems {
				mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashResult)
				preloadedCache.Put(item.key, item.value)
			}
		})

		Context("Given an empty cache", func() {
			It("should return an empty map", func() {
				emptyCache, err := NewCache[int, any](5, LRU_ALGO)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(emptyCache.ListAll())).Should(BeZero())
			})
		})

		Context("Given a preloaded cache with 4 elements", func() {
			It("should return a map with preloaded elements", func() {
				allItems := preloadedCache.ListAll()

				Expect(len(allItems)).Should(Equal(4))
				// checking all preloaded items
				for _, item := range preloadedItems {
					v, ok := allItems[item.key]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(BeEquivalentTo(item.value))
				}
			})
		})
	})
}

// deleteTest tests Delete functionality
func deleteTest() {
	When("Testing DELETE function", func() {
		var (
			mockedHashKeyToIntConverter *hashKeyToIntConverterMock[int]
			preloadedCache              *Cache[int, any]
			preloadedItems              = []itemsToLoadWithMockedKeyHashed{
				{
					key:             123,
					value:           "foo",
					mockedHashedKey: 0,
				},
				{
					key:             456,
					value:           "bar",
					mockedHashedKey: 1,
				},
				{
					key:             789,
					value:           struct{}{},
					mockedHashedKey: 0,
				},
				{
					key:             100,
					value:           true,
					mockedHashedKey: 0,
				},
			}
		)

		BeforeEach(func() {
			mockedHashKeyToIntConverter = new(hashKeyToIntConverterMock[int])
			preloadedCache = &Cache[int, any]{
				setSize:               4,
				sets:                  make(map[int]*list.List),
				entries:               make(map[int]*list.Element),
				hashKeyToIntConverter: mockedHashKeyToIntConverter,
				getItemToRemove:       LRU_ITEM_TO_REMOVE_GETTER,
			}

			for _, item := range preloadedItems {
				mockedHashKeyToIntConverter.On("hashKeyToInt", item.key).Return(item.mockedHashedKey)
				preloadedCache.Put(item.key, item.value)
			}
		})

		Context("Trying to delete an item that doesn't exist", func() {
			It("should keep the elements size per key", func() {
				preloadedCache.Delete(999)

				allItems := preloadedCache.ListAll()
				Expect(len(allItems)).Should(Equal(4))
				// checking all preloaded items
				for _, item := range preloadedItems {
					v, ok := allItems[item.key]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(BeEquivalentTo(item.value))
				}
			})
		})

		Context("Trying to delete two items the first one exists and the another one no", func() {
			It("should reduce all times size in 1", func() {
				preloadedCache.Delete(999)
				preloadedCache.Delete(123)

				allItems := preloadedCache.ListAll()
				Expect(len(allItems)).Should(Equal(3))
				// checking all preloaded items
				for k, v := range allItems {
					found := false
					for _, currentItem := range preloadedItems {
						if k == currentItem.key {
							found = true
							Expect(k).Should(BeEquivalentTo(currentItem.key))
							Expect(v).Should(BeEquivalentTo(currentItem.value))
							break
						}
					}
					Expect(found).Should(BeTrue())
				}
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

// concurrencyTest tests thread safe
func concurrencyTest() {
	When("Creating goroutines manually", func() {
		It("Should not return race condition errors if run test with -race flag", func() {
			cache, err := NewCache[string, any](8)
			Expect(err).ShouldNot(HaveOccurred())

			var wg sync.WaitGroup
			wg.Add(6)

			go func() {
				defer wg.Done()
				cache.Put("key1", 100)
			}()

			go func() {
				defer wg.Done()
				cache.Put("key2", true)
			}()

			go func() {
				defer wg.Done()
				cache.ListAll()
			}()

			go func() {
				defer wg.Done()
				cache.ListAll()
			}()

			go func() {
				defer wg.Done()
				cache.Delete("foo")
			}()

			go func() {
				defer wg.Done()
				cache.Get("foo")
			}()

			wg.Wait()
		})
	})

	When("Creating goroutines in a iterative way", func() {
		It("Should not return race condition errors if run test with -race flag", func() {
			cache, err := NewCache[string, any](5)
			Expect(err).ShouldNot(HaveOccurred())

			var wg sync.WaitGroup
			numGoroutines := 100

			wg.Add(numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func(i int) {
					defer wg.Done()
					cache.Put("key", i)     // Concurrently writing
					_, _ = cache.Get("key") // Concurrently reading
				}(i)
			}

			wg.Wait() // Wait for all goroutines to finish
		})
	})
}
