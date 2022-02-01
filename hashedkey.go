package embedded

import (
	"constraints"
	"fmt"
	"hash/fnv"
)

type HashMapKeyType interface {
	constraints.Ordered
}

type hashKey[TKey HashMapKeyType] struct {
	value TKey
	hash  int
}

func newHashKey[TKey HashMapKeyType](key TKey) hashKey[TKey] {
	h := fnv.New64()
	_, _ = h.Write([]byte(fmt.Sprint(key)))
	return hashKey[TKey]{
		value: key,
		hash:  int(h.Sum64()),
	}
}
