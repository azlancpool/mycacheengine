package cacheiface

type Cache[K comparable, V any] interface {
	Put(key K, value V)
	Get(key K) (V, bool)
}
