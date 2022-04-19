package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type listEntry struct {
	data int
	link embedded.ListLink[listValue]
}

var listEntryLinkField = unsafe.Offsetof(listEntry{}.link)

const listDefaultSize = 1000

type (
	listValue     listEntry
	listType      embedded.List[listValue]
	listSetupFunc func(size int) listType
)

func (h listValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func listSetup(size int) listType {
	return embedded.NewList[listValue](listEntryLinkField)
}

func TestEmbeddedList(t *testing.T) {
	list := listSetup(listDefaultSize)
	data := make([]listValue, listDefaultSize)
	for i := 0; i < len(data); i++ {
		data[i].data = i
	}
	t.Run("InsertLast", func(t *testing.T) {
		for i := 2; i < len(data)-1; i++ {
			expected := &data[i]
			if result := list.InsertLast(expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		}
	})
	t.Run("InsertFirst", func(t *testing.T) {
		expected := &data[0]
		if result := list.InsertFirst(expected); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("InsertBefore", func(t *testing.T) {
		expected := &data[1]
		if result := list.InsertBefore(&data[2], expected); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("InsertAfter", func(t *testing.T) {
		expected := &data[len(data)-1]
		if result := list.InsertAfter(&data[len(data)-2], expected); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("First", func(t *testing.T) {
		expected := &data[0]
		if result := list.First(); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Next", func(t *testing.T) {
		entry := list.First()
		expected := &data[1]
		if result := list.Next(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Last", func(t *testing.T) {
		expected := &data[len(data)-1]
		if result := list.Last(); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Prev", func(t *testing.T) {
		entry := list.Last()
		expected := &data[len(data)-2]
		if result := list.Prev(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("IsContained", func(t *testing.T) {
		t.Run("Contained", func(t *testing.T) {
			expected := true
			if result := list.IsContained(&data[0]); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Uncontained", func(t *testing.T) {
			expected := false
			entry := &listValue{data: 0}
			if result := list.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Nil", func(t *testing.T) {
			expected := false
			var entry *listValue
			if result := list.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
	})
	t.Run("RemoveAll", func(t *testing.T) {
		list.RemoveAll()
	})
}

func BenchmarkEmbeddedList(b *testing.B) {
	list := listSetup(listDefaultSize)
	data := make([]listValue, b.N)
	for i := 0; i < len(data); i++ {
		data[i].data = i
	}

	b.Run("InsertLast", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := range data {
			list.InsertLast(&data[i])
		}
	})
	b.Run("First", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = list.First()
		}
	})
	b.Run("Next", func(b *testing.B) {
		entry := list.First()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entry = list.Next(entry)
		}
	})
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = list.Last()
		}
	})
	b.Run("Prev", func(b *testing.B) {
		entry := list.Last()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entry = list.Prev(entry)
		}
	})
	b.Run("IsContained", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = list.IsContained(&data[i%len(data)])
		}
	})
	b.Run("RemoveAll", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		list.RemoveAll()
	})
}
