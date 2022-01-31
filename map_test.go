package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type mapEntry struct {
	data int
	link embedded.MapLink[int, mapEntry]
}

var mapEntryLinkField = unsafe.Offsetof(mapEntry{}.link)

func TestEmbeddedMap(t *testing.T) {
	const testSize = 5500
	m := embedded.NewMap[int, mapEntry](mapEntryLinkField)
	for i := 0; i < testSize; i++ {
		m.Insert(i, &mapEntry{data: i})
	}

	cur := m.Last()
	for i := testSize - 1; i >= 0; i-- {
		if cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
		cur = m.Prev(cur)
	}
}

func BenchmarkEmbeddedMap1k(b *testing.B) {
	size := 1000
	m := embedded.NewMap[int, mapEntry](mapEntryLinkField)
	for i := 0; i < size; i++ {
		m.Insert(i, &mapEntry{data: i})
	}
}

func BenchmarkEmbeddedMap10k(b *testing.B) {
	size := 10000
	m := embedded.NewMap[int, mapEntry](mapEntryLinkField)
	for i := 0; i < size; i++ {
		m.Insert(i, &mapEntry{data: i})
	}
}

func BenchmarkEmbeddedMap50k(b *testing.B) {
	size := 50000
	m := embedded.NewMap[int, mapEntry](mapEntryLinkField)
	for i := 0; i < size; i++ {
		m.Insert(i, &mapEntry{data: i})
	}
}
