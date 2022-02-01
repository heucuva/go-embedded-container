package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashMapEntry struct {
	data int
	link embedded.HashMapLink[int, hashMapEntry]
}

var hashMapEntryLinkField = unsafe.Offsetof(hashMapEntry{}.link)

func TestEmbeddedHashMapStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	hashMap := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, staticSize)
	for i := 0; i < testSize; i++ {
		hashMap.Insert(i, &hashMapEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hashMap.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func TestEmbeddedHashMapDynamic(t *testing.T) {
	const testSize = 5500
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < testSize; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hash.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func BenchmarkEmbeddedHashMapStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}
