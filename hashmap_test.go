package embedded_test

import (
	"math/rand"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashMapEntry struct {
	data int
	link embedded.HashMapLink[hashMapKey, hashMapValue]
}

var hashMapEntryLinkField = unsafe.Offsetof(hashMapEntry{}.link)

const hashMapDefaultSize = 1000

type (
	hashMapKey       int
	hashMapValue     hashMapEntry
	hashMapType      embedded.HashMap[hashMapKey, hashMapValue]
	hashMapSetupFunc func(size int) hashMapType
)

func (h hashMapValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func hashMapSetupStatic(size int) hashMapType {
	return embedded.NewHashMapStatic[hashMapKey, hashMapValue](hashMapEntryLinkField, size)
}

func hashMapSetupDynamic(size int) hashMapType {
	return embedded.NewHashMapDynamic[hashMapKey, hashMapValue](hashMapEntryLinkField)
}

func TestEmbeddedHashMap(t *testing.T) {
	t.Run("Static", hashMapTest(hashMapSetupStatic))
	t.Run("Dynamic", hashMapTest(hashMapSetupDynamic))
}

func hashMapTest(setupFunc hashMapSetupFunc) func(t *testing.T) {
	return func(t *testing.T) {
		hashMap := setupFunc(hashMapDefaultSize)
		data := make([]hashMapValue, hashMapDefaultSize)
		for i := 0; i < len(data); i++ {
			key := hashMapKey(i)
			data[key].data = i
		}
		t.Run("Reserve", func(t *testing.T) {
			if res := hashMapTestReserve(hashMap, hashMapDefaultSize*1.75); res != nil {
				if !hashMap.IsStatic() {
					t.Fatal("dynamic hashMap is expected to successfully reserve")
				}
			}
		})
		t.Run("Insert", func(t *testing.T) {
			for i := 0; i < len(data); i++ {
				key := hashMapKey(i)
				expected := &data[key]
				if result := hashMap.Insert(key, expected); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindFirst", func(t *testing.T) {
			for i := range data {
				key := hashMapKey(i)
				expected := &data[key]
				if result := hashMap.FindFirst(key); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindNext", func(t *testing.T) {
			for i := range data {
				key := hashMapKey(i)
				entry := &data[key]
				var expected *hashMapValue
				if result := hashMap.FindNext(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("IsContained", func(t *testing.T) {
			t.Run("Contained", func(t *testing.T) {
				expected := true
				key := hashMapKey(0)
				if result := hashMap.IsContained(&data[key]); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("Uncontained", func(t *testing.T) {
				expected := false
				entry := &hashMapValue{data: 0}
				if result := hashMap.IsContained(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("Nil", func(t *testing.T) {
				expected := false
				var entry *hashMapValue
				if result := hashMap.IsContained(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("RemoveAll", func(t *testing.T) {
			hashMap.RemoveAll()
		})
	}
}

func hashMapTestReserve(hashMap hashMapType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	hashMap.Reserve(size)
	err = nil
	return
}

func BenchmarkEmbeddedHashMap(b *testing.B) {
	b.Run("Static", hashMapBench(hashMapSetupStatic))
	b.Run("Dynamic", hashMapBench(hashMapSetupDynamic))
}

func hashMapBench(setupFunc hashMapSetupFunc) func(b *testing.B) {
	return func(b *testing.B) {
		hashMap := setupFunc(hashMapDefaultSize)
		data := make([]hashMapValue, b.N)
		for i := 0; i < len(data); i++ {
			key := hashMapKey(i)
			data[key].data = i
		}

		b.Run("IsStatic", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashMap.IsStatic()
		})
		b.Run("Reserve", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			_ = hashMapBenchReserve(hashMap, int(float64(b.N)*1.75))
		})
		b.Run("Insert", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := range data {
				key := hashMapKey(i)
				hashMap.Insert(key, &data[key])
			}
		})
		b.Run("FindFirst", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := hashMapKey(i % len(data))
				_ = hashMap.FindFirst(key)
			}
		})
		b.Run("FindNext", func(b *testing.B) {
			key := hashMapKey(int(rand.Int31()) % len(data))
			entry := hashMap.FindFirst(key)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hashMap.FindNext(entry)
			}
		})
		b.Run("IsContained", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := hashMapKey(i % len(data))
				_ = hashMap.IsContained(&data[key])
			}
		})
		b.Run("RemoveAll", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hashMap.RemoveAll()
		})
	}
}

func hashMapBenchReserve(hashMap hashMapType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	if size > 1000000 {
		// too big
		err = "too big"
		return
	}
	hashMap.Reserve(size)
	err = nil
	return
}
