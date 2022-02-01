package embedded

import (
	"unsafe"
)

// PriorityQueueLink is a link to the priority queue container
type PriorityQueueLink[P PriorityType] struct {
	position int
	priority P
}

func getPriorityQueueLink[P PriorityType, T any](obj *T, linkFieldOfs uintptr) *PriorityQueueLink[P] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*PriorityQueueLink[P])(u)
}
