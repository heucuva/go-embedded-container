package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type listEntry struct {
	data int
	link embedded.ListLink[listEntry]
}

var listEntryLinkField = unsafe.Offsetof(listEntry{}.link)

func TestEmbeddedList(t *testing.T) {
	const testSize = 5500
	list := embedded.NewList[listEntry](listEntryLinkField)
	for i := 0; i < testSize; i++ {
		list.InsertLast(&listEntry{data: i})
	}

	cur := list.Last()
	for i := testSize - 1; i >= 0; i-- {
		if cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
		cur = list.Prev(cur)
	}
}

func BenchmarkEmbeddedList1k(b *testing.B) {
	size := 1000
	list := embedded.NewList[listEntry](listEntryLinkField)
	for i := 0; i < size; i++ {
		list.InsertLast(&listEntry{data: i})
	}
}

func BenchmarkEmbeddedList100k(b *testing.B) {
	size := 100000
	list := embedded.NewList[listEntry](listEntryLinkField)
	for i := 0; i < size; i++ {
		list.InsertLast(&listEntry{data: i})
	}
}

func BenchmarkEmbeddedList1M(b *testing.B) {
	size := 1000000
	list := embedded.NewList[listEntry](listEntryLinkField)
	for i := 0; i < size; i++ {
		list.InsertLast(&listEntry{data: i})
	}
}
