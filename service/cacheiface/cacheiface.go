package cacheiface

type Cache[K comparable, V any] interface {
	Put(key K, value V)
	Get(key K) (V, bool)
	ListAll() map[K]V
	Delete(key K)
}
