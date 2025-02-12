package cache

import (
	"fmt"
	"hash/fnv"
)

// hashKeyToInt converts a given key of any comparable type into a hashed integer value.
func hashKeyToInt[K comparable](key K) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%v", key)))
	return int(hasher.Sum32())
}
