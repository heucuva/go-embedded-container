package embedded

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

func NewHashStatic[T any](linkField uintptr, tableSize int) Hash[T] {
	return &embeddedHashStatic[T]{
		linkField: linkField,
		table:     make([]*T, tableSize),
	}
}

type embeddedHashStatic[T any] struct {
	entryCount int
	linkField  uintptr
	table      []*T
}

func (c *embeddedHashStatic[T]) Insert(hashValue int, obj *T) *T {
	spot := int(uint(hashValue) % uint(len(c.table)))
	u := getHashLink(obj, c.linkField)
	u.hashValue = hashValue
	u.hashNext = c.table[spot]
	c.table[spot] = obj
	c.entryCount++
	return obj
}

func (c *embeddedHashStatic[T]) Remove(obj *T) *T {
	spot := int(uint(getHashLink(obj, c.linkField).hashValue) % uint(len(c.table)))
	cur := c.table[spot]
	prev := &c.table[spot]

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

func (c *embeddedHashStatic[T]) Move(obj *T, newHashValue int) {
	c.Remove(obj)
	c.Insert(newHashValue, obj)
}

func (c *embeddedHashStatic[T]) Reserve(count int) {
	panic("cannot reserve a count with a static table size")
}

func (c *embeddedHashStatic[T]) GetKey(obj *T) int {
	return getHashLink(obj, c.linkField).hashValue
}

func (c *embeddedHashStatic[T]) Count() int {
	return c.entryCount
}

func (c *embeddedHashStatic[T]) GetTableSize() int {
	return len(c.table)
}

func (c *embeddedHashStatic[T]) GetTableUsed() int {
	if c.entryCount <= 1 {
		return c.entryCount
	}

	var tableUsed int
	for _, entry := range c.table {
		if entry != nil {
			tableUsed++
		}
	}
	return tableUsed
}

func (c *embeddedHashStatic[T]) IsEmpty() bool {
	return c.entryCount == 0
}

func (c *embeddedHashStatic[T]) FindFirst(hashValue int) *T {
	spot := int(uint(hashValue) % uint(len(c.table)))
	entry := c.table[spot]
	for entry != nil {
		u := getHashLink(entry, c.linkField)
		if u.hashValue == hashValue {
			return entry
		}
		entry = u.hashNext
	}
	return nil
}

func (c *embeddedHashStatic[T]) FindNext(prevResult *T) *T {
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

func (c *embeddedHashStatic[T]) WalkFirst() *T {
	if c.entryCount == 0 {
		return nil
	}

	for _, entry := range c.table {
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHashStatic[T]) WalkNext(prevResult *T) *T {
	entry := prevResult
	u := getHashLink(entry, c.linkField)
	spot := int(uint(u.hashValue) % uint(len(c.table)))
	entry = u.hashNext
	if entry != nil {
		return entry
	}

	for spot++; spot < len(c.table); spot++ {
		entry = c.table[spot]
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHashStatic[T]) RemoveAll() {
	cur := c.WalkFirst()
	for cur != nil {
		next := c.WalkNext(cur)
		c.Remove(cur)
		cur = next
	}
}

func (c *embeddedHashStatic[T]) IsContained(cur *T) bool {
	u := getHashLink(cur, c.linkField)
	if u.hashValue != 0 || u.hashNext != nil {
		return true
	}

	walk := c.table[0]
	for walk != nil {
		if walk == cur {
			return true
		}
		walk = getHashLink(walk, c.linkField).hashNext
	}
	return false
}
