package embedded

import (
	"github.com/heucuva/go-embedded-container/internal/array"
)

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type Hash[T any] interface {
	Insert(hashValue int, obj *T) *T
	Remove(obj *T) *T
	Move(obj *T, newHashValue int)
	Reserve(count int)
	GetKey(obj *T) int
	Count() int
	GetTableSize() int
	GetTableUsed() int
	IsEmpty() bool
	FindFirst(hashValue int) *T
	FindNext(prevResult *T) *T
	WalkFirst() *T
	WalkNext(prevResult *T) *T
	RemoveAll()
	IsContained(cur *T) bool
}

func NewHashStatic[T any](linkField uintptr, tableSize int) Hash[T] {
	return &embeddedHash[T]{
		linkField: linkField,
		table:     array.NewStaticArray[*T](tableSize),
	}
}

func NewHashDynamic[T any](linkField uintptr) Hash[T] {
	h := &embeddedHash[T]{
		linkField: linkField,
	}
	h.table = array.NewDynamicArray[*T](minDynamicHashSize, h.onResize)
	return h
}

const (
	minDynamicHashSize = 8
)

type embeddedHash[T any] struct {
	entryCount int
	linkField  uintptr
	table      array.Array[*T]
}

func (c *embeddedHash[T]) Insert(hashValue int, obj *T) *T {
	if !c.table.IsStatic() {
		c.Reserve(c.entryCount + 1)
	}
	spot := int(uint(hashValue) % uint(c.table.Size()))
	u := getHashLink(obj, c.linkField)
	u.hashValue = hashValue
	u.hashNext = c.table.Slice()[spot]
	c.table.Slice()[spot] = obj
	c.entryCount++
	return obj
}

func (c *embeddedHash[T]) Remove(obj *T) *T {
	spot := int(uint(getHashLink(obj, c.linkField).hashValue) % uint(c.table.Size()))
	cur := c.table.Slice()[spot]
	prev := &c.table.Slice()[spot]

	for cur != nil {
		u := getHashLink(cur, c.linkField)
		if cur == obj {
			*prev = u.hashNext
			u.hashNext = nil
			u.hashValue = 0
			c.entryCount--
			return cur
		}
		prev = &u.hashNext
		cur = u.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) Move(obj *T, newHashValue int) {
	c.Remove(obj)
	c.Insert(newHashValue, obj)
}

func (c *embeddedHash[T]) Reserve(count int) {
	if c.table.IsStatic() {
		panic("cannot reserve a count with a static table size")
	} else {
		c.table.Reserve(count)
	}
}

func (c *embeddedHash[T]) GetKey(obj *T) int {
	return getHashLink(obj, c.linkField).hashValue
}

func (c *embeddedHash[T]) Count() int {
	return c.entryCount
}

func (c *embeddedHash[T]) GetTableSize() int {
	return c.table.Size()
}

func (c *embeddedHash[T]) GetTableUsed() int {
	if c.entryCount <= 1 {
		return c.entryCount
	}

	var tableUsed int
	for _, entry := range c.table.Slice() {
		if entry != nil {
			tableUsed++
		}
	}
	return tableUsed
}

func (c *embeddedHash[T]) IsEmpty() bool {
	return c.entryCount == 0
}

func (c *embeddedHash[T]) FindFirst(hashValue int) *T {
	spot := int(uint(hashValue) % uint(c.table.Size()))
	entry := c.table.Slice()[spot]
	for entry != nil {
		u := getHashLink(entry, c.linkField)
		if u.hashValue == hashValue {
			return entry
		}
		entry = u.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) FindNext(prevResult *T) *T {
	entry := prevResult
	u := getHashLink(entry, c.linkField)
	hashValue := u.hashValue
	entry = u.hashNext
	for entry != nil {
		u = getHashLink(entry, c.linkField)
		if u.hashValue == hashValue {
			return entry
		}
		entry = u.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) WalkFirst() *T {
	if c.entryCount == 0 {
		return nil
	}

	for _, entry := range c.table.Slice() {
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHash[T]) WalkNext(prevResult *T) *T {
	entry := prevResult
	u := getHashLink(entry, c.linkField)
	spot := int(uint(u.hashValue) % uint(c.table.Size()))
	entry = u.hashNext
	if entry != nil {
		return entry
	}

	for spot++; spot < c.table.Size(); spot++ {
		entry = c.table.Slice()[spot]
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHash[T]) RemoveAll() {
	cur := c.WalkFirst()
	for cur != nil {
		next := c.WalkNext(cur)
		c.Remove(cur)
		cur = next
	}
}

func (c *embeddedHash[T]) IsContained(cur *T) bool {
	u := getHashLink(cur, c.linkField)
	if u.hashValue != 0 || u.hashNext != nil {
		return true
	}

	walk := c.table.Slice()[0]
	for walk != nil {
		if walk == cur {
			return true
		}
		walk = getHashLink(walk, c.linkField).hashNext
	}
	return false
}

func (c *embeddedHash[T]) onResize(dest, src []*T) {
	if c.entryCount == 0 {
		return
	}

	dynamicSize := uint(len(dest))
	for _, current := range src {
		if current == nil {
			continue
		}

		var tempBucketRoot *T
		for current != nil {
			u := getHashLink(current, c.linkField)
			next := u.hashNext
			u.hashNext = tempBucketRoot
			tempBucketRoot = current
			current = next
		}

		current = tempBucketRoot
		for current != nil {
			u := getHashLink(current, c.linkField)
			next := u.hashNext
			spot := int(uint(u.hashValue) % dynamicSize)
			u.hashNext = dest[spot]
			dest[spot] = current
			current = next
		}
	}
}
