package cache

import (
	"container/list"
	"fmt"
	"hash/fnv"
)

type Cache[K comparable, V any] struct {
	setSize               int
	sets                  map[int]*list.List
	entries               map[K]*list.Element
	hashKeyToIntConverter hashKeyToIntConverter[K]
	getItemToRemove       func(currentSet *list.List) *list.Element
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

type ReplacementAlgo string

const (
	LRU_ALGO ReplacementAlgo = "LRU"
	MRU_ALGO ReplacementAlgo = "MRU"
)

var (
	LRU_ITEM_TO_REMOVE_GETTER = func(currentSet *list.List) *list.Element {
		return currentSet.Back()
	}

	MRU_ITEM_TO_REMOVE_GETTER = func(currentSet *list.List) *list.Element {
		return currentSet.Front()
	}
)

// NewCache returns a new instance of Cache. It saves the provided setSize in the returned instance
func NewCache[K comparable, V any](setSize int, replacementAlgorithm ReplacementAlgo) (*Cache[K, V], error) {
	if setSize <= 0 {
		return nil, fmt.Errorf("setSize provided '%d', must be a positive value", setSize)
	}

	var zero K
	if !isPrimitiveDataType(zero) {
		return nil, fmt.Errorf("provided data type for key is not a supported primitive data type, data type received: %T", zero)
	}

	getItemToRemove := LRU_ITEM_TO_REMOVE_GETTER
	if replacementAlgorithm == MRU_ALGO {
		getItemToRemove = MRU_ITEM_TO_REMOVE_GETTER
	}

	return &Cache[K, V]{
		setSize:               setSize,
		sets:                  make(map[int]*list.List),
		entries:               make(map[K]*list.Element),
		hashKeyToIntConverter: new(hashKeyToIntImpl[K]),
		getItemToRemove:       getItemToRemove,
	}, nil
}

// Put implements functionality that seet a new value in the cache, following n-way-set-associative-cache
func (c *Cache[K, V]) Put(key K, value V) {
	setIndex := c.hashKeyToIntConverter.hashKeyToInt(key) % c.setSize
	if elem, found := c.entries[key]; found {
		c.sets[setIndex].MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	if c.sets[setIndex] == nil {
		c.sets[setIndex] = list.New()
	}

	if c.sets[setIndex].Len() >= c.setSize {
		elementToRemove := c.getItemToRemove(c.sets[setIndex])
		if elementToRemove != nil {
			delete(c.entries, elementToRemove.Value.(*entry[K, V]).key)
			c.sets[setIndex].Remove(elementToRemove)
		}
	}

	newEntry := &entry[K, V]{key: key, value: value}
	elem := c.sets[setIndex].PushFront(newEntry)
	c.entries[key] = elem
}

// Get returns the item if it's present in cache and a true flag.
// Otherwise it returns false and an empty value
func (c *Cache[K, V]) Get(key K) (V, bool) {
	if elem, found := c.entries[key]; found {
		setIndex := c.hashKeyToIntConverter.hashKeyToInt(key) % c.setSize
		c.sets[setIndex].MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// isPrimitiveDataType returns true if the input data type is int, float32, float64, bool or string
func isPrimitiveDataType[K any](input K) bool {
	switch any(input).(type) {
	case int, float32, float64, bool, string:
		return true
	}

	return false
}

type hashKeyToIntConverter[K comparable] interface {
	hashKeyToInt(key K) int
}

type hashKeyToIntImpl[K comparable] struct{}

// hashKeyToInt converts a given key of any comparable type into a hashed integer value.
// For this it get the string equivalent to concatenation of the key value and data type.
func (*hashKeyToIntImpl[K]) hashKeyToInt(key K) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%[1]v%[1]T", key)))
	return int(hasher.Sum32())
}
