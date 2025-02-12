package cache

import (
	"container/list"
	"fmt"
	"hash/fnv"
)

type Cache[K comparable, V any] struct {
	setSize int
	sets    map[int]*list.List
	entries map[K]*list.Element
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewCache returns a new instance of Cache. It saves the provided setSize in the returned instance
func NewCache[K comparable, V any](setSize int) (*Cache[K, V], error) {
	if setSize <= 0 {
		return nil, fmt.Errorf("setSize provided '%d', must be a positive value", setSize)
	}

	var zero K
	if !isPrimitiveDataType(zero) {
		return nil, fmt.Errorf("provided data type for key is not a supported primitive data type, data type received: %T", zero)
	}

	return &Cache[K, V]{
		setSize: setSize,
		sets:    make(map[int]*list.List),
		entries: make(map[K]*list.Element),
	}, nil
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
