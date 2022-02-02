package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashMapEntry struct {
	data int
	link embedded.HashMapLink[int, hashMapEntry]
}

var hashMapEntryLinkField = unsafe.Offsetof(hashMapEntry{}.link)

func TestEmbeddedHashMapStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	const expectedTableUsed = 996
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, staticSize)
	testEmbeddedHashMap(t, c, testSize, expectedTableUsed, staticSize, removeTarget)
}

func TestEmbeddedHashMapDynamic(t *testing.T) {
	const testSize = 5500
	const expectedTableUsed = 4116
	const expectedTableSize = 8192 // next power of 2 over 5500
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	testEmbeddedHashMap(t, c, testSize, expectedTableUsed, expectedTableSize, removeTarget)
}

func BenchmarkEmbeddedHashMapStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashMapStatic[int, hashMapEntry](hashMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashMapDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashMapDynamic[int, hashMapEntry](hashMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.Insert(i, &hashMapEntry{data: i})
	}
}

func testEmbeddedHashMap(t *testing.T, c embedded.HashMap[int, hashMapEntry], testSize, expectedTableUsed, expectedTableSize int, removeTarget int) {
	for i := 0; i < testSize; i++ {
		c.Insert(i, &hashMapEntry{data: i})
	}

	if c.IsEmpty() {
		t.Fatal("embedded hash should not be empty")
	}

	if actualTableSize := c.GetTableSize(); actualTableSize != expectedTableSize {
		t.Fatalf("unexpected table size (actual %d != expected %d)", actualTableSize, expectedTableSize)
	}

	if actualTableUsed := c.GetTableUsed(); actualTableUsed != expectedTableUsed {
		t.Fatalf("unexpected table used size (actual %d != expected %d)", actualTableUsed, expectedTableUsed)
	}

	var removedEntry *hashMapEntry
	for i := testSize - 1; i >= 0; i-- {
		var entry *hashMapEntry
		for cur := c.FindFirst(i); cur != nil; cur = c.FindNext(cur) {
			if cur.data == i {
				entry = cur
				break
			}
		}
		if entry == nil {
			t.Fatal("expected entry not found")
		}

		if actualKey := c.GetKey(entry); actualKey != i {
			t.Fatalf("hashed key mismatch detected (actual %d != expected %d", actualKey, i)
		}

		if i == removeTarget {
			if !c.IsContained(entry) {
				t.Fatal("embedded hash reports that contained item is not present")
			}
			removedEntry = c.Remove(entry)
			if c.IsContained(entry) {
				t.Fatal("embedded hash reports that removed item is present")
			}
		}
	}

	expectedCount := c.Count()

	if moveItem := c.WalkFirst(); moveItem != nil {
		oldKey := c.GetKey(moveItem)
		newKey := testSize
		for i := 1; newKey == oldKey; i++ {
			newKey = testSize + i
		}
		c.Move(moveItem, newKey)
		currentKey := c.GetKey(moveItem)
		if currentKey != newKey {
			t.Fatalf("moved item did not move to expected key hash (old %d -> actual %d != expected %d", oldKey, currentKey, newKey)
		}
	} else {
		t.Fatal("could not find any item in embedded hash")
	}

	if actualCount := c.Count(); actualCount != expectedCount {
		t.Fatalf("count changed unexpectedly (actual %d != expected %d)", actualCount, expectedCount)
	}

	c.Insert(removeTarget, removedEntry)

	c.RemoveAll()
	if actualTableUsed := c.GetTableUsed(); actualTableUsed != 0 {
		t.Fatalf("unexpected table used size (actual %d != expected %d)", actualTableUsed, 0)
	}
}
