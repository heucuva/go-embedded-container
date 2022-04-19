package embedded_test

import (
	"math/rand"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListMapEntry struct {
	data int
	link embedded.HashListMapLink[hashListMapKey, hashListMapValue]
}

var hashListMapEntryLinkField = unsafe.Offsetof(hashListMapEntry{}.link)

const hashListMapDefaultSize = 1000

type (
	hashListMapKey       int
	hashListMapValue     hashListMapEntry
	hashListMapType      embedded.HashListMap[hashListMapKey, hashListMapValue]
	hashListMapSetupFunc func(size int) hashListMapType
)

func (h hashListMapValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func hashListMapSetupStatic(size int) hashListMapType {
	return embedded.NewHashListMapStatic[hashListMapKey, hashListMapValue](hashListMapEntryLinkField, size)
}

func hashListMapSetupDynamic(size int) hashListMapType {
	return embedded.NewHashListMapDynamic[hashListMapKey, hashListMapValue](hashListMapEntryLinkField)
}

func TestEmbeddedHashListMap(t *testing.T) {
	t.Run("Static", hashListMapTest(hashListMapSetupStatic))
	t.Run("Dynamic", hashListMapTest(hashListMapSetupDynamic))
}

func hashListMapTest(setupFunc hashListMapSetupFunc) func(t *testing.T) {
	return func(t *testing.T) {
		hashListMap := setupFunc(hashListMapDefaultSize)
		data := make([]hashListMapValue, hashListMapDefaultSize)
		for i := 0; i < len(data); i++ {
			key := hashListMapKey(i)
			data[key].data = i
		}
		t.Run("Reserve", func(t *testing.T) {
			if res := hashListMapTestReserve(hashListMap, hashListMapDefaultSize*1.75); res != nil {
				if !hashListMap.IsStatic() {
					t.Fatal("dynamic hashListMap is expected to successfully reserve")
				}
			}
		})
		t.Run("InsertLast", func(t *testing.T) {
			for i := 2; i < len(data)-2; i++ {
				key := hashListMapKey(i)
				expected := &data[key]
				if result := hashListMap.InsertLast(key, expected); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("InsertFirst", func(t *testing.T) {
			key := hashListMapKey(0)
			expected := &data[key]
			if result := hashListMap.InsertFirst(key, expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("InsertBefore", func(t *testing.T) {
			key := hashListMapKey(1)
			expected := &data[key]
			if result := hashListMap.InsertBefore(key, &data[key+1], expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("InsertAfter", func(t *testing.T) {
			key := hashListMapKey(len(data) - 1)
			expected := &data[key]
			if result := hashListMap.InsertAfter(key, &data[key-1], expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("FindFirst", func(t *testing.T) {
			for i := range data {
				key := hashListMapKey(i)
				expected := &data[key]
				if result := hashListMap.FindFirst(key); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindNext", func(t *testing.T) {
			for i := range data {
				key := hashListMapKey(i)
				entry := &data[key]
				var expected *hashListMapValue
				if result := hashListMap.FindNext(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("First", func(t *testing.T) {
			key := hashListMapKey(0)
			expected := &data[key]
			if result := hashListMap.First(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Next", func(t *testing.T) {
			entry := hashListMap.First()
			key := hashListMapKey(1)
			expected := &data[key]
			if result := hashListMap.Next(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Last", func(t *testing.T) {
			key := hashListMapKey(len(data) - 1)
			expected := &data[key]
			if result := hashListMap.Last(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Prev", func(t *testing.T) {
			key := hashListMapKey(len(data) - 2)
			entry := hashListMap.Last()
			expected := &data[key]
			if result := hashListMap.Prev(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("IsContained", func(t *testing.T) {
			key := hashListMapKey(0)
			expected := true
			if result := hashListMap.IsContained(&data[key]); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
			expected = false
			var entry *hashListMapValue
			if result := hashListMap.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
			entry = &hashListMapValue{data: 0}
			if result := hashListMap.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("RemoveAll", func(t *testing.T) {
			hashListMap.RemoveAll()
		})
	}
}

func hashListMapTestReserve(hashListMap hashListMapType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	hashListMap.Reserve(size)
	err = nil
	return
}

func BenchmarkEmbeddedHashListMap(b *testing.B) {
	b.Run("Static", hashListMapBench(hashListMapSetupStatic))
	b.Run("Dynamic", hashListMapBench(hashListMapSetupDynamic))
}

func hashListMapBench(setupFunc hashListMapSetupFunc) func(b *testing.B) {
	return func(b *testing.B) {
		hashListMap := setupFunc(hashListMapDefaultSize)
		data := make([]hashListMapValue, b.N)
		for i := 0; i < len(data); i++ {
			key := hashListMapKey(i)
			data[key].data = i
		}

		b.Run("IsStatic", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashListMap.IsStatic()
		})
		b.Run("Reserve", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			_ = hashListMapBenchReserve(hashListMap, int(float64(b.N)*1.75))
		})
		b.Run("InsertLast", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := range data {
				key := hashListMapKey(i)
				hashListMap.InsertLast(key, &data[key])
			}
		})
		b.Run("FindFirst", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := hashListMapKey(i % len(data))
				_ = hashListMap.FindFirst(key)
			}
		})
		b.Run("FindNext", func(b *testing.B) {
			key := hashListMapKey(int(rand.Int31()) % len(data))
			entry := hashListMap.FindFirst(key)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashListMap.FindNext(entry)
			}
		})
		b.Run("First", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashListMap.First()
			}
		})
		b.Run("Next", func(b *testing.B) {
			entry := hashListMap.First()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashListMap.Next(entry)
			}
		})
		b.Run("Last", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashListMap.Last()
			}
		})
		b.Run("Prev", func(b *testing.B) {
			entry := hashListMap.Last()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashListMap.Prev(entry)
			}
		})
		b.Run("IsContained", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := hashListMapKey(i % len(data))
				_ = hashListMap.IsContained(&data[key])
			}
		})
		b.Run("RemoveAll", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashListMap.RemoveAll()
		})
	}
}

func hashListMapBenchReserve(hashListMap hashListMapType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	if size > 1000000 {
		// too big
		err = "too big"
		return
	}
	hashListMap.Reserve(size)
	err = nil
	return
}
