package embedded_test

import (
	"math/rand"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashEntry struct {
	data int
	link embedded.HashLink[hashValue]
}

var hashEntryLinkField = unsafe.Offsetof(hashEntry{}.link)

const hashDefaultSize = 1000

type (
	hashValue     hashEntry
	hashType      embedded.Hash[hashValue]
	hashSetupFunc func(size int) hashType
)

func (h hashValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func hashSetupStatic(size int) hashType {
	return embedded.NewHashStatic[hashValue](hashEntryLinkField, size)
}

func hashSetupDynamic(size int) hashType {
	return embedded.NewHashDynamic[hashValue](hashEntryLinkField)
}

func TestEmbeddedHash(t *testing.T) {
	t.Run("Static", hashTest(hashSetupStatic))
	t.Run("Dynamic", hashTest(hashSetupDynamic))
}

func hashTest(setupFunc hashSetupFunc) func(t *testing.T) {
	return func(t *testing.T) {
		hash := setupFunc(hashDefaultSize)
		data := make([]hashValue, hashDefaultSize)
		for i := 0; i < len(data); i++ {
			data[i].data = i
		}
		t.Run("Reserve", func(t *testing.T) {
			if res := hashTestReserve(hash, hashDefaultSize*1.75); res != nil {
				if !hash.IsStatic() {
					t.Fatal("dynamic hash is expected to successfully reserve")
				}
			}
		})
		t.Run("Insert", func(t *testing.T) {
			for i := range data {
				expected := &data[i]
				if result := hash.Insert(expected.Hash(), expected); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindFirst", func(t *testing.T) {
			for i := range data {
				expected := &data[i]
				if result := hash.FindFirst(expected.Hash()); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindNext", func(t *testing.T) {
			for i := range data {
				entry := &data[i]
				var expected *hashValue
				if result := hash.FindNext(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("WalkFirst", func(t *testing.T) {
			expected := &data[0]
			if result := hash.WalkFirst(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("WalkNext", func(t *testing.T) {
			entry := hash.WalkFirst()
			expected := &data[1]
			if result := hash.WalkNext(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("IsContained", func(t *testing.T) {
			t.Run("Contained", func(t *testing.T) {
				expected := true
				if result := hash.IsContained(&data[0]); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("Uncontained", func(t *testing.T) {
				expected := false
				entry := &hashValue{data: 0}
				if result := hash.IsContained(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("Nil", func(t *testing.T) {
				expected := false
				var entry *hashValue
				if result := hash.IsContained(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("RemoveAll", func(t *testing.T) {
			hash.RemoveAll()
		})
	}
}

func hashTestReserve(hash hashType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	hash.Reserve(size)
	err = nil
	return
}

func BenchmarkEmbeddedHash(b *testing.B) {
	b.Run("Static", hashBench(hashSetupStatic))
	b.Run("Dynamic", hashBench(hashSetupDynamic))
}

func hashBench(setupFunc hashSetupFunc) func(b *testing.B) {
	return func(b *testing.B) {
		hash := setupFunc(hashDefaultSize)
		data := make([]hashValue, b.N)
		for i := 0; i < len(data); i++ {
			data[i].data = i
		}

		b.Run("IsStatic", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hash.IsStatic()
		})
		b.Run("Reserve", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			_ = hashBenchReserve(hash, int(float64(b.N)*1.75))
		})
		b.Run("Insert", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := range data {
				hash.Insert(data[i].Hash(), &data[i])
			}
		})
		b.Run("FindFirst", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hash.FindFirst(data[i%len(data)].Hash())
			}
		})
		b.Run("FindNext", func(b *testing.B) {
			entry := hash.FindFirst(data[int(rand.Int31())%len(data)].Hash())
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hash.FindNext(entry)
			}
		})
		b.Run("WalkFirst", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hash.WalkFirst()
			}
		})
		b.Run("WalkNext", func(b *testing.B) {
			entry := hash.WalkFirst()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				entry = hash.WalkNext(entry)
			}
		})
		b.Run("IsContained", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hash.IsContained(&data[i%len(data)])
			}
		})
		b.Run("RemoveAll", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hash.RemoveAll()
		})
	}
}

func hashBenchReserve(hash hashType, size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	if size > 1000000 {
		// too big
		err = "too big"
		return
	}
	hash.Reserve(size)
	err = nil
	return
}
