package embedded_test

import (
	"math/rand"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListEntry struct {
	data int
	link embedded.HashListLink[hashListValue]
}

var hashListEntryLinkField = unsafe.Offsetof(hashListEntry{}.link)

const hashListDefaultSize = 1000

type (
	hashListValue     hashListEntry
	hashListType      embedded.HashList[hashListValue]
	hashListSetupFunc func(size int) hashListType
)

func (h hashListValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func hashListSetupStatic(size int) hashListType {
	return embedded.NewHashListStatic[hashListValue](hashListEntryLinkField, size)
}

func hashListSetupDynamic(size int) hashListType {
	return embedded.NewHashListDynamic[hashListValue](hashListEntryLinkField)
}

func TestEmbeddedHashList(t *testing.T) {
	t.Run("Static", hashListTest(hashListSetupStatic))
	t.Run("Dynamic", hashListTest(hashListSetupDynamic))
}

func hashListTest(setupFunc hashListSetupFunc) func(t *testing.T) {
	return func(t *testing.T) {
		hashList := setupFunc(hashListDefaultSize)
		data := make([]hashListValue, hashListDefaultSize)
		for i := 0; i < len(data); i++ {
			data[i].data = i
		}
		t.Run("Reserve", func(t *testing.T) {
			if res := hashListTestReserve(hashList, hashListDefaultSize*1.75); res != nil {
				if !hashList.IsStatic() {
					t.Fatal("dynamic hashList is expected to successfully reserve")
				}
			}
		})
		t.Run("InsertLast", func(t *testing.T) {
			for i := 2; i < len(data)-2; i++ {
				expected := &data[i]
				if result := hashList.InsertLast(expected.Hash(), expected); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("InsertFirst", func(t *testing.T) {
			expected := &data[0]
			if result := hashList.InsertFirst(expected.Hash(), expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("InsertBefore", func(t *testing.T) {
			expected := &data[1]
			if result := hashList.InsertBefore(expected.Hash(), &data[2], expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("InsertAfter", func(t *testing.T) {
			expected := &data[len(data)-1]
			if result := hashList.InsertAfter(expected.Hash(), &data[len(data)-2], expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("FindFirst", func(t *testing.T) {
			for i := range data {
				expected := &data[i]
				if result := hashList.FindFirst(expected.Hash()); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindNext", func(t *testing.T) {
			for i := range data {
				entry := &data[i]
				var expected *hashListValue
				if result := hashList.FindNext(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("First", func(t *testing.T) {
			expected := &data[0]
			if result := hashList.First(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Next", func(t *testing.T) {
			entry := hashList.First()
			expected := &data[1]
			if result := hashList.Next(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Last", func(t *testing.T) {
			expected := &data[len(data)-1]
			if result := hashList.Last(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Prev", func(t *testing.T) {
			entry := hashList.Last()
			expected := &data[len(data)-2]
			if result := hashList.Prev(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("IsContained", func(t *testing.T) {
			expected := true
			if result := hashList.IsContained(&data[0]); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
			expected = false
			var entry *hashListValue
			if result := hashList.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
			entry = &hashListValue{data: 0}
			if result := hashList.IsContained(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("RemoveAll", func(t *testing.T) {
			hashList.RemoveAll()
		})
	}
}

func hashListTestReserve(hashList hashListType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	hashList.Reserve(size)
	err = nil
	return
}

func BenchmarkEmbeddedHashList(b *testing.B) {
	b.Run("Static", hashListBench(hashListSetupStatic))
	b.Run("Dynamic", hashListBench(hashListSetupDynamic))
}

func hashListBench(setupFunc hashListSetupFunc) func(b *testing.B) {
	return func(b *testing.B) {
		hashList := setupFunc(hashListDefaultSize)
		data := make([]hashListValue, b.N)
		for i := 0; i < len(data); i++ {
			data[i].data = i
		}

		b.Run("IsStatic", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashList.IsStatic()
		})
		b.Run("Reserve", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			_ = hashListBenchReserve(hashList, int(float64(b.N)*1.75))
		})
		b.Run("InsertLast", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := range data {
				hashList.InsertLast(data[i].Hash(), &data[i])
			}
		})
		b.Run("FindFirst", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashList.FindFirst(data[i%len(data)].Hash())
			}
		})
		b.Run("FindNext", func(b *testing.B) {
			entry := hashList.FindFirst(data[int(rand.Int31())%len(data)].Hash())
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashList.FindNext(entry)
			}
		})
		b.Run("First", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashList.First()
			}
		})
		b.Run("Next", func(b *testing.B) {
			entry := hashList.First()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashList.Next(entry)
			}
		})
		b.Run("Last", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashList.Last()
			}
		})
		b.Run("Prev", func(b *testing.B) {
			entry := hashList.Last()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashList.Prev(entry)
			}
		})
		b.Run("IsContained", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hashList.IsContained(&data[i%len(data)])
			}
		})
		b.Run("RemoveAll", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashList.RemoveAll()
		})
	}
}

func hashListBenchReserve(hashList hashListType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	if size > 1000000 {
		// too big
		err = "too big"
		return
	}
	hashList.Reserve(size)
	err = nil
	return
}
