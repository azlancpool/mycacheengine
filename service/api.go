package cache

import (
	"fmt"
	"hash/fnv"
)

// hashKeyToInt converts a given key of any comparable type into a hashed integer value.
// For this it get the string equivalent to concatenation of the key value and data type.
func hashKeyToInt[K comparable](key K) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%[1]v%[1]T", key)))
	return int(hasher.Sum32())
}
