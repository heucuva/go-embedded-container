package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type hashEntry struct {
	data int
	link embedded.HashLink[hashEntry]
}

var hashEntryLinkField = unsafe.Offsetof(hashEntry{}.link)

func TestEmbeddedHashStatic(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	const expectedTableUsed = 996
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashStatic[hashEntry](hashEntryLinkField, staticSize)
	testEmbeddedHash(t, c, testSize, expectedTableUsed, staticSize, removeTarget)
}

func TestEmbeddedHashStaticReserve(t *testing.T) {
	const staticSize = 1000
	const testSize = int(staticSize * 5.5)
	const expectedTableUsed = 996
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashStatic[hashEntry](hashEntryLinkField, staticSize)
	defer func() {
		err := recover()
		if err != nil {
			t.SkipNow()
		}
	}()
	c.Reserve(staticSize)
	t.FailNow()
}

func TestEmbeddedHashDynamic(t *testing.T) {
	const testSize = 5500
	const expectedTableUsed = 4116
	const expectedTableSize = 8192 // next power of 2 over 5500
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	testEmbeddedHash(t, c, testSize, expectedTableUsed, expectedTableSize, removeTarget)
}

func BenchmarkEmbeddedHashStatic_Insert(b *testing.B) {
	hash := embedded.NewHashStatic[hashEntry](hashEntryLinkField, b.N)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hkey := embedded.HashKey(i)
		hash.Insert(hkey, &hashEntry{data: i})
	}
}

func BenchmarkEmbeddedHashDynamic_Insert(b *testing.B) {
	hash := embedded.NewHashDynamic[hashEntry](hashEntryLinkField)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hkey := embedded.HashKey(i)
		hash.Insert(hkey, &hashEntry{data: i})
	}
}

func testEmbeddedHash(t *testing.T, c embedded.Hash[hashEntry], testSize, expectedTableUsed, expectedTableSize int, removeTarget int) {
	for i := 0; i < testSize; i++ {
		hkey := embedded.HashKey(i)
		c.Insert(hkey, &hashEntry{data: i})
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

	for i := testSize - 1; i >= 0; i-- {
		hkey := embedded.HashKey(i)
		var entry *hashEntry
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
			if !c.IsContained(entry) {
				t.Fatal("embedded hash reports that contained item is not present")
			}
			c.Remove(entry)
			if c.IsContained(entry) {
				t.Fatal("embedded hash reports that removed item is present")
			}
		}
	}

	for walk := c.WalkFirst(); walk != nil; walk = c.WalkNext(walk) {
		if walk.data == removeTarget {
			t.Fatal("removed item still present in embedded hash")
		}
	}

	expectedCount := c.Count()

	if moveItem := c.WalkFirst(); moveItem != nil {
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

	c.RemoveAll()
	if actualTableUsed := c.GetTableUsed(); actualTableUsed != 0 {
		t.Fatalf("unexpected table used size (actual %d != expected %d)", actualTableUsed, 0)
	}
}
