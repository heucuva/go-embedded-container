package embedded

import (
	"encoding/binary"
	"hash/fnv"
	"reflect"

	"golang.org/x/exp/constraints"
)

type HashMapKeyType interface {
	constraints.Ordered
}

type HashedKeyValue uint64

type Hashable interface {
	Hash() HashedKeyValue
}

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
	if v, ok := (any(key)).(Hashable); ok {
		return v.Hash()
	} else {
		v := reflect.ValueOf(key)
		if v.CanUint() {
			return HashedKeyValue(v.Uint())
		} else if v.CanInt() {
			return HashedKeyValue(uint64(v.Int()))
		}
	}

	h := fnv.New64()
	binary.Write(h, binary.BigEndian, key)
	return HashedKeyValue(h.Sum64())
}
