package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListMapEntry struct {
	data int
	link embedded.HashListMapLink[int, hashListMapEntry]
}

var hashListMapEntryLinkField = unsafe.Offsetof(hashListMapEntry{}.link)

func TestEmbeddedHashListMapStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	const expectedTableUsed = 996
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, staticSize)
	testEmbeddedHashListMap(t, c, testSize, expectedTableUsed, staticSize, removeTarget)
}

func TestEmbeddedHashListMapDynamic(t *testing.T) {
	const testSize = 5500
	const expectedTableUsed = 4116
	const expectedTableSize = 8192 // next power of 2 over 5500
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	testEmbeddedHashListMap(t, c, testSize, expectedTableUsed, expectedTableSize, removeTarget)
}

func BenchmarkEmbeddedHashListMapStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListMapStatic[int, hashListMapEntry](hashListMapEntryLinkField, size)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListMapDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListMapDynamic[int, hashListMapEntry](hashListMapEntryLinkField)
	for i := 0; i < size; i++ {
		hash.InsertLast(i, &hashListMapEntry{data: i})
	}
}

func testEmbeddedHashListMap(t *testing.T, c embedded.HashListMap[int, hashListMapEntry], testSize, expectedTableUsed, expectedTableSize int, removeTarget int) {
	for i := 0; i < testSize; i++ {
		c.InsertLast(i, &hashListMapEntry{data: i})
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

	var removedEntry *hashListMapEntry
	for i := testSize - 1; i >= 0; i-- {
		var entry *hashListMapEntry
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
			t.Fatalf("hashed key mismatch detected (actual %d != expected %d)", actualKey, i)
		}

		if i == removeTarget {
			if actualEntry := c.Position(i); actualEntry != entry {
				t.Fatalf("item not found at expected position")
			}
			if !c.IsContained(entry) {
				t.Fatal("embedded hash reports that contained item is not present")
			}
			removedEntry = c.Remove(entry)
			if c.IsContained(entry) {
				t.Fatal("embedded hash reports that removed item is present")
			}
		}
	}

	for walk := c.First(); walk != nil; walk = c.Next(walk) {
		if walk.data == removeTarget {
			t.Fatal("removed item still present in embedded hash")
		}
	}

	for walk := c.Last(); walk != nil; walk = c.Prev(walk) {
		if walk.data == removeTarget {
			t.Fatal("removed item still present in embedded hash")
		}
	}

	expectedCount := c.Count()

	if moveItem := c.First(); moveItem != nil {
		oldKey := c.GetKey(moveItem)
		newKey := testSize
		for i := 1; newKey == oldKey; i++ {
			newKey = testSize + i
		}
		c.Move(moveItem, newKey)
		currentKey := c.GetKey(moveItem)
		if currentKey != newKey {
			t.Fatalf("moved item did not move to expected key hash (old %d -> actual %d != expected %d)", oldKey, currentKey, newKey)
		}
	} else {
		t.Fatal("could not find any item in embedded hash")
	}

	if actualCount := c.Count(); actualCount != expectedCount {
		t.Fatalf("count changed unexpectedly (actual %d != expected %d)", actualCount, expectedCount)
	}

	c.InsertFirst(removeTarget, removedEntry)

	if actualFirst := c.RemoveFirst(); actualFirst == nil {
		t.Fatal("no item at front of embedded hash list")
	} else if expectedFirst := removeTarget; actualFirst.data != expectedFirst {
		t.Fatalf("mismatched item at front of embedded hash list (actual %d != expected %d)", actualFirst.data, expectedFirst)
	}

	if actualLast := c.RemoveLast(); actualLast == nil {
		t.Fatal("no item at front of embedded hash list")
	} else if expectedLast := testSize - 1; actualLast.data != expectedLast {
		t.Fatalf("mismatched item at front of embedded hash list (actual %d != expected %d)", actualLast.data, expectedLast)
	}

	c.RemoveAll()
	if actualTableUsed := c.GetTableUsed(); actualTableUsed != 0 {
		t.Fatalf("unexpected table used size (actual %d != expected %d)", actualTableUsed, 0)
	}
}
