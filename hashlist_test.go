package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListEntry struct {
	data int
	link embedded.HashListLink[hashListEntry]
}

var hashListEntryLinkField = unsafe.Offsetof(hashListEntry{}.link)

func TestEmbeddedHashListStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	hashList := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, staticSize)
	for i := 0; i < testSize; i++ {
		hashList.InsertLast(i, &hashListEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hashList.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func TestEmbeddedHashListDynamic(t *testing.T) {
	const testSize = 5500
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < testSize; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}

	for i := testSize - 1; i >= 0; i-- {
		if cur := hash.FindFirst(i); cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func BenchmarkEmbeddedHashListStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListEntry{data: i})
	}
}
