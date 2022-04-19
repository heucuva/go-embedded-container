package embedded_test

import (
	"math/rand"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type mapEntry struct {
	data int
	link embedded.MapLink[mapKey, mapValue]
}

var mapEntryLinkField = unsafe.Offsetof(mapEntry{}.link)

const mapDefaultSize = 1000

type (
	mapKey       int
	mapValue     mapEntry
	mapType      embedded.Map[mapKey, mapValue]
	mapSetupFunc func(size int) mapType
)

func (h mapValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func mapSetup(size int) mapType {
	return embedded.NewMap[mapKey, mapValue](mapEntryLinkField)
}

func TestEmbeddedMap(t *testing.T) {
	m := mapSetup(mapDefaultSize)
	data := make([]mapValue, mapDefaultSize)
	for i := 0; i < len(data); i++ {
		key := mapKey(i)
		data[key].data = i
	}
	t.Run("Insert", func(t *testing.T) {
		for i := 0; i < len(data); i++ {
			key := mapKey(i)
			expected := &data[key]
			if result := m.Insert(key, expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		}
	})
	t.Run("FindFirst", func(t *testing.T) {
		for i := range data {
			key := mapKey(i)
			expected := &data[key]
			if result := m.FindFirst(key); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		}
	})
	t.Run("FindNext", func(t *testing.T) {
		for i := range data {
			key := mapKey(i)
			entry := &data[key]
			var expected *mapValue
			if result := m.FindNext(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		}
	})
	t.Run("First", func(t *testing.T) {
		key := mapKey(0)
		expected := &data[key]
		if result := m.First(); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Next", func(t *testing.T) {
		entry := m.First()
		key := mapKey(1)
		expected := &data[key]
		if result := m.Next(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Last", func(t *testing.T) {
		key := mapKey(len(data) - 1)
		expected := &data[key]
		if result := m.Last(); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("Prev", func(t *testing.T) {
		key := mapKey(len(data) - 2)
		entry := m.Last()
		expected := &data[key]
		if result := m.Prev(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("IsContained", func(t *testing.T) {
		key := mapKey(0)
		expected := true
		if result := m.IsContained(&data[key]); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
		expected = false
		var entry *mapValue
		if result := m.IsContained(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
		entry = &mapValue{data: 0}
		if result := m.IsContained(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("RemoveAll", func(t *testing.T) {
		m.RemoveAll()
	})
}

func BenchmarkEmbeddedMap(b *testing.B) {
	m := mapSetup(mapDefaultSize)
	data := make([]mapValue, mapDefaultSize)
	for i := 0; i < len(data); i++ {
		key := mapKey(i)
		data[key].data = i
	}

	b.Run("Insert", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := range data {
			key := mapKey(i)
			m.Insert(key, &data[key])
		}
	})
	b.Run("FindFirst", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := mapKey(i % len(data))
			_ = m.FindFirst(key)
		}
	})
	b.Run("FindNext", func(b *testing.B) {
		key := mapKey(int(rand.Int31()) % len(data))
		entry := m.FindFirst(key)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entry = m.FindNext(entry)
		}
	})
	b.Run("First", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m.First()
		}
	})
	b.Run("Next", func(b *testing.B) {
		entry := m.First()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entry = m.Next(entry)
		}
	})
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = m.Last()
		}
	})
	b.Run("Prev", func(b *testing.B) {
		entry := m.Last()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			entry = m.Prev(entry)
		}
	})
	b.Run("IsContained", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := mapKey(i % len(data))
			_ = m.IsContained(&data[key])
		}
	})
	b.Run("RemoveAll", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		m.RemoveAll()
	})
}
