package embedded

import (
	"constraints"
	"unsafe"
)

// This is a priority queue container - it allows for prioritization of its
// contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type PriorityQueue[P PriorityType, T any] interface {
	Top() *T
	TopWithPriority(priority P) *T
	Insert(priority P, entry *T) *T
	Remove(entry *T) *T
	RemoveTop() *T
	RemoveTopWithPriority(priority P) *T
	RemoveAll()
	Count() int
	IsEmpty() bool

	GetPriority(entry *T) *P
	IsContained(entry *T) bool
}

type PriorityType interface{
	constraints.Integer | constraints.Float
}

func NewPriorityQueue[P PriorityType, T any](linkField uintptr) PriorityQueue[P, T] {
	return &embeddedPriorityQueue[P, T]{
		linkField: linkField,
	}
}

type embeddedPriorityQueue[P PriorityType, T any] struct {
	array     []*T
	linkField uintptr
}

func (c *embeddedPriorityQueue[P, T]) Top() *T {
	if len(c.array) == 0 {
		return nil
	}
	return c.array[0]
}

func (c *embeddedPriorityQueue[P, T]) TopWithPriority(priority P) *T {
	if len(c.array) == 0 {
		return nil
	}
	top := c.array[0]
	u := getPriorityQueueLink[P](top, c.linkField)
	if !(priority < u.priority) {
		return top
	}
	return nil
}

func (c *embeddedPriorityQueue[P, T]) Remove(entry *T) *T {
	if len(c.array) == 0 {
		return entry
	}

	u := getPriorityQueueLink[P](entry, c.linkField)
	spot := int(u.position) - 1
	if spot == -1 {
		return entry
	}

	endEntry := c.array[len(c.array)-1]
	c.array = c.array[:len(c.array)-1]
	if entry != endEntry {
		c.array[spot] = endEntry
		v := getPriorityQueueLink[P](endEntry, c.linkField)
		v.position = spot + 1
		c.refloat(endEntry)
	}
	u.position = 0
	return entry
}

func (c *embeddedPriorityQueue[P, T]) RemoveTop() *T {
	return c.Remove(c.Top())
}

func (c *embeddedPriorityQueue[P, T]) RemoveTopWithPriority(priority P) *T {
	top := c.Top()
	if top == nil {
		return nil
	}

	u := getPriorityQueueLink[P](top, c.linkField)
	if !(priority < u.priority) {
		return c.Remove(top)
	}

	return nil
}

func (c *embeddedPriorityQueue[P, T]) Count() int {
	return len(c.array)
}

func (c *embeddedPriorityQueue[P, T]) IsEmpty() bool {
	return len(c.array) == 0
}

func (c *embeddedPriorityQueue[P, T]) RemoveAll() {
	for len(c.array) > 0 {
		last := c.array[len(c.array)-1]
		c.Remove(last)
	}
}

func (c *embeddedPriorityQueue[P, T]) IsContained(entry *T) bool {
	return c.GetPriority(entry) != nil
}

func (c *embeddedPriorityQueue[P, T]) GetPriority(entry *T) *P {
	u := getPriorityQueueLink[P](entry, c.linkField)
	spot := int(u.priority) - 1
	if spot >= 0 {
		return &u.priority
	}
	return nil
}

func (c *embeddedPriorityQueue[P, T]) Insert(priority P, entry *T) *T {
	u := getPriorityQueueLink[P](entry, c.linkField)
	spot := int(u.position) - 1
	if spot == -1 {
		u.priority = priority
		u.position = len(c.array) + 1
		c.array = append(c.array, entry)
	} else {
		if !(u.priority < priority || priority < u.priority) {
			return entry
		}
		u.priority = priority
	}
	c.refloat(entry)
	return entry
}

func (c *embeddedPriorityQueue[P, T]) refloat(entry *T) {
	u := getPriorityQueueLink[P](entry, c.linkField)
	spot := int(u.position) - 1
	tryDown := true
	for spot > 0 {
		hold := c.array[spot]
		v := getPriorityQueueLink[P](hold, c.linkField)
		newSpot := (spot - 1) / 2
		lower := c.array[newSpot]
		w := getPriorityQueueLink[P](lower, c.linkField)
		if !(v.priority < w.priority) {
			break
		}

		c.array[spot] = c.array[newSpot]
		c.array[newSpot] = hold
		w.position = spot + 1
		v.position = newSpot + 1
		spot = newSpot
		tryDown = false
	}

	if tryDown {
		for {
			downSpot1 := (spot * 2) + 1
			if downSpot1 >= len(c.array) {
				break
			}

			u := getPriorityQueueLink[P](c.array[spot], c.linkField)
			v := getPriorityQueueLink[P](c.array[downSpot1], c.linkField)

			downSpot2 := (spot * 2) + 2
			var w *PriorityQueueLink[P]
			if downSpot2 < len(c.array) {
				w = getPriorityQueueLink[P](c.array[downSpot2], c.linkField)
			}
			if w == nil || v.priority < w.priority {
				if !(v.priority < u.priority) {
					break
				}

				c.array[spot], c.array[downSpot1] = c.array[downSpot1], c.array[spot]
				u.position, v.position = downSpot1+1, spot+1
				spot = downSpot1
			} else {
				if !(w.priority < u.priority) {
					break
				}

				c.array[spot], c.array[downSpot2] = c.array[downSpot2], c.array[spot]
				u.position, w.position = downSpot2+1, spot+1
				spot = downSpot2
			}
		}
	}
}

// PriorityQueueLink is a link to the priority queue container
type PriorityQueueLink[P PriorityType] struct {
	position int
	priority P
}

func getPriorityQueueLink[P PriorityType, T any](obj *T, linkFieldOfs uintptr) *PriorityQueueLink[P] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*PriorityQueueLink[P])(u)
}
