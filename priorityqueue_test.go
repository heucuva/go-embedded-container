package embedded_test

import (
	"testing"
	"unsafe"

	embedded "github.com/heucuva/go-embedded-container"
)

type priorityQueueEntry struct {
	data int
	link embedded.PriorityQueueLink[priorityQueuePriority]
}

var priorityQueueEntryLinkField = unsafe.Offsetof(priorityQueueEntry{}.link)

const priorityQueueDefaultSize = 1000

type (
	priorityQueuePriority  int
	priorityQueueValue     priorityQueueEntry
	priorityQueueType      embedded.PriorityQueue[priorityQueuePriority, priorityQueueValue]
	priorityQueueSetupFunc func(size int) priorityQueueType
)

func (h priorityQueueValue) Hash() embedded.HashedKeyValue {
	return embedded.HashedKeyValue(h.data)
}

func priorityQueueSetup(size int) priorityQueueType {
	return embedded.NewPriorityQueue[priorityQueuePriority, priorityQueueValue](priorityQueueEntryLinkField)
}

func TestEmbeddedPriorityQueue(t *testing.T) {
	priorityQueue := priorityQueueSetup(priorityQueueDefaultSize)
	data := make([]priorityQueueValue, priorityQueueDefaultSize)
	for i := 0; i < len(data); i++ {
		data[i].data = i
	}
	t.Run("Insert", func(t *testing.T) {
		for i := 0; i < len(data); i++ {
			priority := priorityQueuePriority(i)
			expected := &data[i]
			if result := priorityQueue.Insert(priority, expected); result != expected {
				t.Fatalf("expected %v, but got %v", expected, result)
			}
		}
	})
	t.Run("Top", func(t *testing.T) {
		expected := &data[0]
		if result := priorityQueue.Top(); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("TopWithPriority", func(t *testing.T) {
		priority := priorityQueuePriority(3)
		expected := &data[3]
		if result := priorityQueue.TopWithPriority(priority); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("IsContained", func(t *testing.T) {
		expected := true
		if result := priorityQueue.IsContained(&data[0]); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
		expected = false
		var entry *priorityQueueValue
		if result := priorityQueue.IsContained(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
		entry = &priorityQueueValue{data: 0}
		if result := priorityQueue.IsContained(entry); result != expected {
			t.Fatalf("expected %v, but got %v", expected, result)
		}
	})
	t.Run("RemoveAll", func(t *testing.T) {
		priorityQueue.RemoveAll()
	})
}

func BenchmarkEmbeddedPriorityQueue(b *testing.B) {
	priorityQueue := priorityQueueSetup(priorityQueueDefaultSize)
	data := make([]priorityQueueValue, b.N)
	for i := 0; i < len(data); i++ {
		data[i].data = i
	}

	b.Run("Insert", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := range data {
			priority := priorityQueuePriority(i)
			priorityQueue.Insert(priority, &data[i])
		}
	})
	b.Run("Top", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = priorityQueue.Top()
		}
	})
	b.Run("TopWithPriority", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			priority := priorityQueuePriority(i % len(data))
			_ = priorityQueue.TopWithPriority(priority)
		}
	})
	b.Run("IsContained", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = priorityQueue.IsContained(&data[i%len(data)])
		}
	})
	b.Run("RemoveAll", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		priorityQueue.RemoveAll()
	})
}
