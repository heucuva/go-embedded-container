package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListMapEntry struct {
	data int
	link embedded.HashListMapLink[int, hashListMapEntry]
}

var hashListMapEntryLinkField = unsafe.Offsetof(hashListMapEntry{}.link)

func TestEmbeddedHashListMapStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	hashListMap := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, staticSize)
	for i := 0; i < testSize; i++ {
		hashListMap.InsertLast(i, &hashListMapEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hashListMap.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func TestEmbeddedHashListMapDynamic(t *testing.T) {
	const testSize = 5500
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < testSize; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hash.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func BenchmarkEmbeddedHashListMapStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}
