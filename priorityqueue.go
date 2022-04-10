package embedded

import (
	"golang.org/x/exp/constraints"
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

type PriorityType interface {
	constraints.Ordered
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

func (c *embeddedPriorityQueue[P, T]) getLink(obj *T) *PriorityQueueLink[P] {
	return getPriorityQueueLink[P](obj, c.linkField)
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
	topLink := c.getLink(top)
	if !(priority < topLink.priority) {
		return top
	}
	return nil
}

func (c *embeddedPriorityQueue[P, T]) Remove(entry *T) *T {
	if len(c.array) == 0 {
		return entry
	}

	entryLink := c.getLink(entry)
	spot := int(entryLink.position) - 1
	if spot == -1 {
		return entry
	}

	endEntry := c.array[len(c.array)-1]
	c.array = c.array[:len(c.array)-1]
	if entry != endEntry {
		c.array[spot] = endEntry
		v := c.getLink(endEntry)
		v.position = spot + 1
		c.refloat(endEntry)
	}
	entryLink.position = 0
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

	topLink := c.getLink(top)
	if !(priority < topLink.priority) {
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
	entryLink := c.getLink(entry)
	spot := int(entryLink.position) - 1
	if spot >= 0 {
		return &entryLink.priority
	}
	return nil
}

func (c *embeddedPriorityQueue[P, T]) Insert(priority P, entry *T) *T {
	entryLink := c.getLink(entry)
	spot := int(entryLink.position) - 1
	if spot == -1 {
		entryLink.priority = priority
		entryLink.position = len(c.array) + 1
		c.array = append(c.array, entry)
	} else {
		if !(entryLink.priority < priority || priority < entryLink.priority) {
			return entry
		}
		entryLink.priority = priority
	}
	c.refloat(entry)
	return entry
}

func (c *embeddedPriorityQueue[P, T]) refloat(entry *T) {
	entryLink := c.getLink(entry)
	spot := int(entryLink.position) - 1
	tryDown := true
	for spot > 0 {
		hold := c.array[spot]
		v := c.getLink(hold)
		newSpot := (spot - 1) / 2
		lower := c.array[newSpot]
		w := c.getLink(lower)
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

			curLink := c.getLink(c.array[spot])
			downSpot1Link := c.getLink(c.array[downSpot1])

			downSpot2 := (spot * 2) + 2
			var downSpot2Link *PriorityQueueLink[P]
			if downSpot2 < len(c.array) {
				downSpot2Link = c.getLink(c.array[downSpot2])
			}
			if downSpot2Link == nil || downSpot1Link.priority < downSpot2Link.priority {
				if !(downSpot1Link.priority < curLink.priority) {
					break
				}

				c.array[spot], c.array[downSpot1] = c.array[downSpot1], c.array[spot]
				curLink.position, downSpot1Link.position = downSpot1+1, spot+1
				spot = downSpot1
			} else {
				if !(downSpot2Link.priority < curLink.priority) {
					break
				}

				c.array[spot], c.array[downSpot2] = c.array[downSpot2], c.array[spot]
				curLink.position, downSpot2Link.position = downSpot2+1, spot+1
				spot = downSpot2
			}
		}
	}
}
