package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashEntry struct {
	data int
	link embedded.HashLink[hashEntry]
}

var hashEntryLinkField = unsafe.Offsetof(hashEntry{}.link)

func TestEmbeddedHashStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	hash := embedded.NewHashStatic[hashEntry](hashEntryLinkField, staticSize)
	for i := 0; i < testSize; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hash.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func TestEmbeddedHashDynamic(t *testing.T) {
	const testSize = 5500
	hash := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	for i := 0; i < testSize; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hash.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func BenchmarkEmbeddedHashStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashStatic[hashEntry](hashEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashStatic[hashEntry](hashEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashStatic[hashEntry](hashEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashEntry{data: i})
	}
}
