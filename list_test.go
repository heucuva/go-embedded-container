package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type listEntry struct {
	data int
	link embedded.ListLink[listEntry]
}

var listEntryLinkField = unsafe.Offsetof(listEntry{}.link)

func TestEmbeddedList(t *testing.T) {
	const testSize = 5500
	const removeTarget = (testSize / 2) - 1
	c := embedded.NewList[listEntry](listEntryLinkField)
	testEmbeddedList(t, c, testSize, removeTarget)
}

func BenchmarkEmbeddedList_InsertLast(b *testing.B) {
	list := embedded.NewList[listEntry](listEntryLinkField)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		list.InsertLast(&listEntry{data: i})
	}
}

func BenchmarkEmbeddedList_InsertFirst(b *testing.B) {
	list := embedded.NewList[listEntry](listEntryLinkField)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		list.InsertFirst(&listEntry{data: i})
	}
}

func testEmbeddedList(t *testing.T, c embedded.List[listEntry], testSize int, removeTarget int) {
	for i := 0; i < testSize; i++ {
		c.InsertLast(&listEntry{data: i})
	}

	if c.IsEmpty() {
		t.Fatal("embedded list should not be empty")
	}

	var removedEntry *listEntry
	for i := testSize - 1; i >= 0; i-- {
		var entry *listEntry
		for cur := c.Last(); cur != nil; cur = c.Prev(cur) {
			if cur.data == i {
				entry = cur
				break
			}
		}
		if entry == nil {
			t.Fatal("expected entry not found")
		}

		if i == removeTarget {
			if actualEntry := c.Position(i); actualEntry != entry {
				t.Fatalf("item not found at expected position")
			}
			if !c.IsContained(entry) {
				t.Fatal("embedded list reports that contained item is not present")
			}
			removedEntry = c.Remove(entry)
			if c.IsContained(entry) {
				t.Fatal("embedded list reports that removed item is present")
			}
		}
	}

	for walk := c.First(); walk != nil; walk = c.Next(walk) {
		if walk.data == removeTarget {
			t.Fatal("removed item still present in embedded list")
		}
	}

	for walk := c.Last(); walk != nil; walk = c.Prev(walk) {
		if walk.data == removeTarget {
			t.Fatal("removed item still present in embedded list")
		}
	}

	expectedCount := c.Count()

	if moveItem := c.First(); moveItem != nil {
		c.MoveLast(moveItem)
		if lastItem := c.Last(); moveItem == nil {
			t.Fatal("could not find moved item in embedded list")
		} else if lastItem != moveItem {
			t.Fatal("expected to find moved item, but found something else")
		}
	} else {
		t.Fatal("could not find any item in embedded list")
	}

	if actualCount := c.Count(); actualCount != expectedCount {
		t.Fatalf("count changed unexpectedly (actual %d != expected %d)", actualCount, expectedCount)
	}

	c.MoveLast(removedEntry)
	c.MoveAfter(c.First(), removedEntry)
	c.MoveBefore(c.Last(), removedEntry)
	c.MoveFirst(removedEntry)

	if actualLast := c.RemoveLast(); actualLast == nil {
		t.Fatal("no item at end of embedded list")
	} else if expectedLast := 0; actualLast.data != expectedLast {
		t.Fatalf("mismatched item at end of embedded list (actual %d != expected %d)", actualLast.data, expectedLast)
	}

	if actualFirst := c.RemoveFirst(); actualFirst == nil {
		t.Fatal("no item at front of embedded list")
	} else if expectedFirst := removeTarget; actualFirst.data != expectedFirst {
		t.Fatalf("mismatched item at front of embedded list (actual %d != expected %d)", actualFirst.data, expectedFirst)
	}

	c.RemoveAll()
	if actualCount := c.Count(); actualCount != 0 {
		t.Fatalf("unexpected list count (actual %d != expected %d)", actualCount, 0)
	}
}
