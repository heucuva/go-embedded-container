package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashListEntry struct {
	data int
	link embedded.HashListLink[hashListEntry]
}

var hashListEntryLinkField = unsafe.Offsetof(hashListEntry{}.link)

func TestEmbeddedHashListStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	const expectedTableUsed = 996
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, staticSize)
	testEmbeddedHashList(t, c, testSize, expectedTableUsed, staticSize, removeTarget)
}

func TestEmbeddedHashListDynamic(t *testing.T) {
	const testSize = 5500
	const expectedTableUsed = 4116
	const expectedTableSize = 8192 // next power of 2 over 5500
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	testEmbeddedHashList(t, c, testSize, expectedTableUsed, expectedTableSize, removeTarget)
}

func BenchmarkEmbeddedHashListStatic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListStatic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListStatic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListStatic[hashListEntry](hashListEntryLinkField, size)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic1k(b *testing.B) {
	size := 1000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic100k(b *testing.B) {
	size := 100000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func BenchmarkEmbeddedHashListDynamic1M(b *testing.B) {
	size := 1000000
	hash := embedded.NewHashListDynamic[hashListEntry](hashListEntryLinkField)
	for i := 0; i < size; i++ {
		hkey := embedded.HashKey(i)
		hash.InsertLast(hkey, &hashListEntry{data: i})
	}
}

func testEmbeddedHashList(t *testing.T, c embedded.HashList[hashListEntry], testSize, expectedTableUsed, expectedTableSize int, removeTarget int) {
	for i := 0; i < testSize; i++ {
		hkey := embedded.HashKey(i)
		c.InsertLast(hkey, &hashListEntry{data: i})
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

	var removedEntry *hashListEntry
	for i := testSize - 1; i >= 0; i-- {
		hkey := embedded.HashKey(i)
		var entry *hashListEntry
		for cur := c.FindFirst(hkey); cur != nil; cur = c.FindNext(cur) {
			if cur.data == i {
				entry = cur
				break
			}
		}
		if entry == nil {
			t.Fatal("expected entry not found")
		}

		if actualHash := c.GetKey(entry); actualHash != hkey {
			t.Fatalf("hashed key mismatch detected (actual %08X != expected %08X", actualHash, hkey)
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
		newKey := embedded.HashKey(testSize)
		for i := 1; newKey == oldKey; i++ {
			newKey = embedded.HashKey(testSize + i)
		}
		c.Move(moveItem, newKey)
		currentKey := c.GetKey(moveItem)
		if currentKey != newKey {
			t.Fatalf("moved item did not move to expected key hash (old %08X -> actual %08X != expected %08X", oldKey, currentKey, newKey)
		}
	} else {
		t.Fatal("could not find any item in embedded hash")
	}

	if actualCount := c.Count(); actualCount != expectedCount {
		t.Fatalf("count changed unexpectedly (actual %d != expected %d)", actualCount, expectedCount)
	}

	removedHkey := embedded.HashKey(removeTarget)
	c.InsertFirst(removedHkey, removedEntry)

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
