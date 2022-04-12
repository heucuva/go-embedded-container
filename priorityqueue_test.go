package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type priorityQueueEntry struct {
	data int
	link embedded.PriorityQueueLink[int]
}

var priorityQueueEntryLinkField = unsafe.Offsetof(priorityQueueEntry{}.link)

func TestEmbeddedPriorityQueue(t *testing.T) {
	const testSize = 5500
	priorityQueue := embedded.NewPriorityQueue[int, priorityQueueEntry](priorityQueueEntryLinkField)
	for i := 0; i < testSize; i++ {
		priorityQueue.Insert(i, &priorityQueueEntry{data: i})
	}

	for i := 0; i < testSize; i++ {
		cur := priorityQueue.RemoveTopWithPriority(i)
		if cur == nil || cur.data != i {
			t.Fatal("expected entry not found")
		}
	}
}

func BenchmarkEmbeddedPriorityQueue_Insert(b *testing.B) {
	priorityQueue := embedded.NewPriorityQueue[int, priorityQueueEntry](priorityQueueEntryLinkField)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		priorityQueue.Insert(i, &priorityQueueEntry{data: i})
	}
}
