package embedded_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
	"github.com/heucuva/go-embedded-container/internal/util"
)

type hashEntry[T any] struct {
	data T
	link embedded.HashLink[hashEntry[T]]
}

const hashDefaultSize = 1000

type (
	hashValueInt        hashEntry[int]
	hashValueString     hashEntry[string]
	hashValueConstraint interface {
		hashValueInt | hashValueString
	}
	hashType[T any]      embedded.Hash[T]
	hashSetupFunc[T any] func(size int) hashType[T]
)

func (h hashValueInt) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func (h *hashValueInt) SetData(i int) {
	h.data = i
}

func (h hashValueString) Hash() embedded.HashedKeyValue {
	var i int
	fmt.Sscan(h.data, &i)
	return embedded.HashKey(i)
}

func (h *hashValueString) SetData(i int) {
	h.data = fmt.Sprint(i)
}

func hashSetupStatic[T any](linkField uintptr, size int) hashType[T] {
	return embedded.NewHashStatic[T](linkField, size)
}

func hashSetupDynamic[T any](linkField uintptr, size int) hashType[T] {
	return embedded.NewHashDynamic[T](linkField)
}

func TestEmbeddedHash(t *testing.T) {
	var hashEntryLinkFieldInt = unsafe.Offsetof(hashValueInt{}.link)
	var hashEntryLinkFieldString = unsafe.Offsetof(hashValueString{}.link)

	t.Run("Static", func(t *testing.T) {
		t.Run("Int", hashTest(func(size int) hashType[hashValueInt] {
			return hashSetupStatic[hashValueInt](hashEntryLinkFieldInt, size)
		}))
		t.Run("String", hashTest(func(size int) hashType[hashValueString] {
			return hashSetupStatic[hashValueString](hashEntryLinkFieldString, size)
		}))
	})
	t.Run("Dynamic", func(t *testing.T) {
		t.Run("Int", hashTest(func(size int) hashType[hashValueInt] {
			return hashSetupDynamic[hashValueInt](hashEntryLinkFieldInt, size)
		}))
		t.Run("String", hashTest(func(size int) hashType[hashValueString] {
			return hashSetupDynamic[hashValueString](hashEntryLinkFieldString, size)
		}))
	})
}

func hashSetData[TValue hashValueConstraint](entry *TValue, i int) {
	switch e := (interface{})(entry).(type) {
	case *hashValueInt:
		e.SetData(i)
	case *hashValueString:
		e.SetData(i)
	default:
		panic(fmt.Errorf("unexpected type %v", reflect.TypeOf(entry)))
	}
}

func hashTest[TValue hashValueConstraint](setupFunc hashSetupFunc[TValue]) func(t *testing.T) {
	return func(t *testing.T) {
		expectedTableSize := hashDefaultSize / 2
		hash := setupFunc(expectedTableSize)
		data := make([]TValue, hashDefaultSize)
		for i := range data {
			entry := &data[i]
			hashSetData(entry, i)
		}
		if !hash.IsStatic() {
			expectedTableSize = int(util.NextPowerOf2(hashDefaultSize + hashDefaultSize>>2))
		}
		t.Run("Insert", func(t *testing.T) {
			for i := range data {
				expected := &data[i]
				h := (interface{})(expected).(embedded.Hashable)
				if result := hash.Insert(h.Hash(), expected); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("Count", func(t *testing.T) {
			expected := len(data)
			if result := hash.Count(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("GetKey", func(t *testing.T) {
			t.Run("Existing", func(t *testing.T) {
				entry := &data[len(data)-1]
				expected := (interface{})(entry).(embedded.Hashable).Hash()
				if result := hash.GetKey(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("NonExisting", func(t *testing.T) {
				var ent TValue
				entry := &ent
				hashSetData(entry, len(data)+2)
				expected := embedded.HashedKeyValue(0)
				if result := hash.GetKey(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("GetTableSize", func(t *testing.T) {
			expected := expectedTableSize
			if result := hash.GetTableSize(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("GetTableUsed", func(t *testing.T) {
			expected := expectedTableSize
			if expectedTableSize > len(data) {
				expected = len(data)
			}
			if result := hash.GetTableUsed(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("IsEmpty", func(t *testing.T) {
			t.Run("Full", func(t *testing.T) {
				expected := false
				if result := hash.IsEmpty(); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("FindFirst", func(t *testing.T) {
			for i := range data {
				expected := &data[i]
				h := (interface{})(expected).(embedded.Hashable)
				if result := hash.FindFirst(h.Hash()); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("FindNext", func(t *testing.T) {
			for i := range data {
				entry := &data[i]
				var expected *TValue
				if result := hash.FindNext(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			}
		})
		t.Run("WalkFirst", func(t *testing.T) {
			i := 0
			if expectedTableSize < len(data) {
				i = expectedTableSize
			}
			expected := &data[i]
			if result := hash.WalkFirst(); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("WalkNext", func(t *testing.T) {
			entry := hash.WalkFirst()
			var expected *TValue
			if expectedTableSize < len(data) {
				expected = &data[0]
			} else {
				expected = &data[1]
			}
			if result := hash.WalkNext(entry); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		})
		t.Run("Reserve", func(t *testing.T) {
			newSize := int(hashDefaultSize * 1.75)
			if res := hashTestReserve(hash, newSize); res != nil {
				if !hash.IsStatic() {
					t.Fatal("dynamic hash is expected to successfully reserve")
				}
			} else if !hash.IsStatic() {
				expectedTableSize = int(util.NextPowerOf2(uint(newSize + newSize>>2)))
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
				var entry TValue
				if result := hash.IsContained(&entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("Nil", func(t *testing.T) {
				expected := false
				var entry TValue
				if result := hash.IsContained(&entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("Remove", func(t *testing.T) {
			t.Run("Existing", func(t *testing.T) {
				entry := &data[len(data)-1]
				expected := entry
				if result := hash.Remove(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
			t.Run("NonExisting", func(t *testing.T) {
				var ent TValue
				entry := &ent
				hashSetData(entry, len(data)+2)
				var expected *TValue
				if result := hash.Remove(entry); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
		t.Run("RemoveAll", func(t *testing.T) {
			hash.RemoveAll()
		})
		t.Run("IsEmpty", func(t *testing.T) {
			t.Run("Empty", func(t *testing.T) {
				expected := true
				if result := hash.IsEmpty(); result != expected {
					t.Fatalf("expected %v, but got %v", expected, result)
				}
			})
		})
	}
}

func hashTestReserve[TValue hashValueConstraint](hash hashType[TValue], size int) (err interface{}) {
	defer func() {
		err = recover()
	}()
	hash.Reserve(size)
	err = nil
	return
}

func BenchmarkEmbeddedHash(b *testing.B) {
	var hashEntryLinkFieldInt = unsafe.Offsetof(hashValueInt{}.link)

	b.Run("Static", hashBench(func(size int) hashType[hashValueInt] {
		return hashSetupStatic[hashValueInt](hashEntryLinkFieldInt, size)
	}))
	b.Run("Dynamic", hashBench(func(size int) hashType[hashValueInt] {
		return hashSetupDynamic[hashValueInt](hashEntryLinkFieldInt, size)
	}))
}

func hashBench(setupFunc hashSetupFunc[hashValueInt]) func(b *testing.B) {
	return func(b *testing.B) {
		hash := setupFunc(hashDefaultSize / 2)
		data := make([]hashValueInt, b.N)
		for i := range data {
			entry := &data[i]
			entry.data = i
		}

		b.Run("IsStatic", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			_ = hash.IsStatic()
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
		b.Run("GetKey", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = hash.GetKey(&data[i%len(data)])
			}
		})
		b.Run("RemoveAll", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			hash.RemoveAll()
		})
	}
}

func hashBenchReserve(hash hashType[hashValueInt], size int) (err interface{}) {
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
