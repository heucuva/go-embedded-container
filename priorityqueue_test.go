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

func BenchmarkEmbeddedPriorityQueue1k(b *testing.B) {
	size := 1000
	priorityQueue := embedded.NewPriorityQueue[int, priorityQueueEntry](priorityQueueEntryLinkField)
	for i := 0; i < size; i++ {
		priorityQueue.Insert(i, &priorityQueueEntry{data: i})
	}
}

func BenchmarkEmbeddedPriorityQueue100k(b *testing.B) {
	size := 100000
	priorityQueue := embedded.NewPriorityQueue[int, priorityQueueEntry](priorityQueueEntryLinkField)
	for i := 0; i < size; i++ {
		priorityQueue.Insert(i, &priorityQueueEntry{data: i})
	}
}

func BenchmarkEmbeddedPriorityQueue1M(b *testing.B) {
	size := 1000000
	priorityQueue := embedded.NewPriorityQueue[int, priorityQueueEntry](priorityQueueEntryLinkField)
	for i := 0; i < size; i++ {
		priorityQueue.Insert(i, &priorityQueueEntry{data: i})
	}
}
