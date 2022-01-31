package embedded

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

func NewHashDynamic[T any](linkField uintptr) Hash[T] {
	return &embeddedHashDynamic[T]{
		linkField: linkField,
	}
}

const (
	minDynamicHashSize = 8
)

type embeddedHashDynamic[T any] struct {
	entryCount int
	linkField  uintptr
	table      []*T
}

func (c *embeddedHashDynamic[T]) Insert(hashValue int, obj *T) *T {
	c.Reserve(c.entryCount + 1)
	spot := int(uint(hashValue) & uint(len(c.table)-1))
	u := getHashLink(obj, c.linkField)
	u.hashValue = hashValue
	u.hashNext = c.table[spot]
	c.table[spot] = obj
	c.entryCount++
	return obj
}

func (c *embeddedHashDynamic[T]) Remove(obj *T) *T {
	spot := int(uint(getHashLink(obj, c.linkField).hashValue) & uint(len(c.table)-1))
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

func (c *embeddedHashDynamic[T]) Move(obj *T, newHashValue int) {
	c.Remove(obj)
	c.Insert(newHashValue, obj)
}

func (c *embeddedHashDynamic[T]) Reserve(count int) {
	count += count >> 2
	if count > len(c.table) {
		c.resize(count)
	}
}

func (c *embeddedHashDynamic[T]) GetKey(obj *T) int {
	return getHashLink(obj, c.linkField).hashValue
}

func (c *embeddedHashDynamic[T]) Count() int {
	return c.entryCount
}

func (c *embeddedHashDynamic[T]) GetTableSize() int {
	return len(c.table)
}

func (c *embeddedHashDynamic[T]) GetTableUsed() int {
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

func (c *embeddedHashDynamic[T]) IsEmpty() bool {
	return c.entryCount == 0
}

func (c *embeddedHashDynamic[T]) FindFirst(hashValue int) *T {
	spot := int(uint(hashValue) & uint(len(c.table)-1))
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

func (c *embeddedHashDynamic[T]) FindNext(prevResult *T) *T {
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

func (c *embeddedHashDynamic[T]) WalkFirst() *T {
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

func (c *embeddedHashDynamic[T]) WalkNext(prevResult *T) *T {
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

func (c *embeddedHashDynamic[T]) RemoveAll() {
	cur := c.WalkFirst()
	for cur != nil {
		next := c.WalkNext(cur)
		c.Remove(cur)
		cur = next
	}
}

func (c *embeddedHashDynamic[T]) IsContained(cur *T) bool {
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

func (c *embeddedHashDynamic[T]) resize(count int) {
	dynamicTableOld := c.table

	dynamicSize := nextPowerOf2(uint(count))
	if dynamicSize < minDynamicHashSize {
		dynamicSize = minDynamicHashSize
	}
	c.table = make([]*T, dynamicSize)
	if c.entryCount == 0 {
		return
	}

	for _, current := range dynamicTableOld {
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
			spot := int(uint(u.hashValue) & uint(dynamicSize-1))
			u.hashNext = c.table[spot]
			c.table[spot] = current
			current = next
		}
	}
}
