package embedded

// This is a double-linked list container - it allows for linear iteration over
// its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type List[T any] interface {
	First() *T
	Last() *T
	Next(cur *T) *T
	Prev(cur *T) *T
	Position(index int) *T

	Remove(obj *T) *T
	RemoveFirst() *T
	RemoveLast() *T
	RemoveAll()

	InsertFirst(cur *T) *T
	InsertLast(cur *T) *T
	InsertAfter(prev, cur *T) *T
	InsertBefore(after, cur *T) *T

	MoveFirst(cur *T)
	MoveLast(cur *T)
	MoveAfter(dest, cur *T)
	MoveBefore(dest, cur *T)

	Count() int
	IsEmpty() bool
	IsContained(cur *T) bool
}

func NewList[T any](linkField uintptr) List[T] {
	return &embeddedList[T]{
		linkField: linkField,
	}
}

type embeddedList[T any] struct {
	head      *T
	tail      *T
	count     int
	linkField uintptr
}

func (c *embeddedList[T]) getLink(obj *T) *ListLink[T] {
	return getListLink(obj, c.linkField)
}

func (c *embeddedList[T]) First() *T {
	return c.head
}

func (c *embeddedList[T]) Last() *T {
	return c.tail
}

func (c *embeddedList[T]) Next(cur *T) *T {
	return c.getLink(cur).next
}

func (c *embeddedList[T]) Prev(cur *T) *T {
	return c.getLink(cur).prev
}

func (c *embeddedList[T]) Position(index int) *T {
	cur := c.head
	for cur != nil && index > 0 {
		cur = c.Next(cur)
		index--
	}
	return cur
}

func (c *embeddedList[T]) Count() int {
	return c.count
}

func (c *embeddedList[T]) Remove(obj *T) *T {
	if c.getLink(obj).remove(c.linkField, &c.head, &c.tail) {
		c.count--
	}
	return obj
}

func (c *embeddedList[T]) RemoveFirst() *T {
	if c.head == nil {
		return nil
	}
	return c.Remove(c.head)
}

func (c *embeddedList[T]) RemoveLast() *T {
	if c.tail == nil {
		return nil
	}
	return c.Remove(c.tail)
}

func (c *embeddedList[T]) RemoveAll() {
	for cur := c.tail; cur != nil; cur = c.tail {
		c.Remove(cur)
	}
}

func (c *embeddedList[T]) InsertFirst(cur *T) *T {
	c.getLink(cur).next = c.head
	if c.head != nil {
		c.getLink(c.head).prev = cur
		c.head = cur
	} else {
		c.head = cur
		c.tail = cur
	}
	c.count++
	return cur
}

func (c *embeddedList[T]) InsertLast(cur *T) *T {
	c.getLink(cur).prev = c.tail
	if c.tail != nil {
		c.getLink(c.tail).next = cur
		c.tail = cur
	} else {
		c.head = cur
		c.tail = cur
	}
	c.count++
	return cur
}

func (c *embeddedList[T]) InsertAfter(prev, cur *T) *T {
	if prev == nil {
		return c.InsertFirst(cur)
	}
	curU := c.getLink(cur)
	prevU := c.getLink(prev)
	curU.prev = prev
	curU.next = prevU.next
	prevU.next = cur

	if curU.next != nil {
		c.getLink(curU.next).prev = cur
	} else {
		c.tail = cur
	}

	c.count++
	return cur
}

func (c *embeddedList[T]) InsertBefore(after, cur *T) *T {
	if after == nil {
		return c.InsertLast(cur)
	}
	curU := c.getLink(cur)
	afterU := c.getLink(after)
	curU.next = after
	curU.prev = afterU.prev
	afterU.prev = cur

	if curU.prev != nil {
		c.getLink(curU.prev).next = cur
	} else {
		c.head = cur
	}

	c.count++
	return cur
}

func (c *embeddedList[T]) MoveFirst(cur *T) {
	c.Remove(cur)
	c.InsertFirst(cur)
}

func (c *embeddedList[T]) MoveLast(cur *T) {
	c.Remove(cur)
	c.InsertLast(cur)
}

func (c *embeddedList[T]) MoveAfter(dest, cur *T) {
	c.Remove(cur)
	c.InsertAfter(dest, cur)
}

func (c *embeddedList[T]) MoveBefore(dest, cur *T) {
	c.Remove(cur)
	c.InsertBefore(dest, cur)
}

func (c *embeddedList[T]) IsEmpty() bool {
	return c.count == 0
}

func (c *embeddedList[T]) IsContained(cur *T) bool {
	return c.getLink(cur).isContained(c.linkField, c.head)
}
