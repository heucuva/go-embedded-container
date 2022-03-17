package embedded

import (
	"fmt"
	"hash/fnv"

	"golang.org/x/exp/constraints"
)

type HashMapKeyType interface {
	constraints.Ordered
}

type HashedKeyValue uint64

type hashKey[TKey HashMapKeyType] struct {
	value TKey
	hash  HashedKeyValue
}

func newHashKey[TKey HashMapKeyType](key TKey) hashKey[TKey] {
	return hashKey[TKey]{
		value: key,
		hash:  HashKey(key),
	}
}

func HashKey[TKey HashMapKeyType](key TKey) HashedKeyValue {
	h := fnv.New64()
	_, _ = h.Write([]byte(fmt.Sprint(key)))
	return HashedKeyValue(h.Sum64())
}
