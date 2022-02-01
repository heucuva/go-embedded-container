package embedded

import (
	"constraints"
)

type Map[TKey MapKeyType, T any] interface {
	Find(key TKey) *T
	FindFirst(key TKey) *T
	FindNext(cur *T) *T
	FindLowerInclusive(key TKey) *T
	FindUpperInclusive(key TKey) *T
	FindLowerExclusive(key TKey) *T
	FindUpperExclusive(key TKey) *T

	First() *T
	Last() *T
	Next(cur *T) *T
	Prev(cur *T) *T
	GetPosition(obj *T) int
	Position(index int) *T

	Insert(key TKey, obj *T) *T
	Remove(obj *T) *T
	RemoveFirst() *T
	RemoveLast() *T
	Move(cur *T, newKey TKey)

	GetKey(obj *T) TKey
	Count() int
	IsEmpty() bool

	RemoveAll()
	IsContained(obj *T) bool
}

type MapKeyType interface {
	constraints.Ordered
}

func NewMap[TKey MapKeyType, T any](linkField uintptr) Map[TKey, T] {
	return &embeddedMap[TKey, T]{
		linkField: linkField,
	}
}

type embeddedMap[TKey MapKeyType, T any] struct {
	root      *T
	count     int
	linkField uintptr
}

func (c *embeddedMap[TKey, T]) Find(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			walk = u.left
		} else if u.key < key {
			walk = u.right
		} else {
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindFirst(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			walk = u.left
		} else if u.key < key {
			walk = u.right
		} else {
			left := c.Prev(walk)
			for left != nil && c.GetKey(left) == key {
				walk = left
				left = c.Prev(walk)
			}
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindNext(cur *T) *T {
	next := c.Next(cur)
	if next != nil && c.GetKey(cur) == c.GetKey(next) {
		return next
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindLowerInclusive(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			left := u.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if u.key < key {
			right := u.right
			if right == nil {
				return walk
			}
			walk = right
		} else {
			left := c.Prev(walk)
			for left != nil && c.GetKey(left) == key {
				walk = left
				left = c.Prev(walk)
			}
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindUpperInclusive(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			left := u.left
			if left == nil {
				return walk
			}
			walk = left
		} else if u.key < key {
			right := u.right
			if right == nil {
				return c.Next(walk)
			}
			walk = right
		} else {
			left := c.Prev(walk)
			for left != nil && c.GetKey(left) == key {
				walk = left
				left = c.Prev(walk)
			}
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindLowerExclusive(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			left := u.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if u.key < key {
			right := u.right
			if right == nil {
				return walk
			}
			walk = right
		} else {
			left := c.Prev(walk)
			for left != nil && c.GetKey(left) == key {
				walk = left
				left = c.Prev(walk)
			}

			if left == nil {
				return nil
			}

			return left
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindUpperExclusive(key TKey) *T {
	walk := c.root
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			left := u.left
			if left == nil {
				return walk
			}
			walk = left
		} else if u.key < key {
			right := u.right
			if right == nil {
				return c.Next(walk)
			}
			walk = right
		} else {
			right := c.Next(walk)
			for right != nil && c.GetKey(right) == key {
				walk = right
				right = c.Next(walk)
			}

			if right == nil {
				return nil
			}

			return right
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) First() *T {
	var prev *T
	cur := c.root
	for cur != nil {
		prev = cur
		cur = getMapLink[TKey](cur, c.linkField).left
	}
	return prev
}

func (c *embeddedMap[TKey, T]) Last() *T {
	var prev *T
	cur := c.root
	for cur != nil {
		prev = cur
		cur = getMapLink[TKey](cur, c.linkField).right
	}
	return prev
}

func (c *embeddedMap[TKey, T]) Next(cur *T) *T {
	u := getMapLink[TKey](cur, c.linkField)
	if u.right != nil {
		walk := u.right
		for {
			v := getMapLink[TKey](walk, c.linkField)
			if v.left == nil {
				break
			}
			walk = v.left
		}
		return walk
	}

	p := u.parent
	for p != nil && getMapLink[TKey](p, c.linkField).right == cur {
		cur = p
		p = getMapLink[TKey](cur, c.linkField).parent
	}
	return p
}

func (c *embeddedMap[TKey, T]) Prev(cur *T) *T {
	u := getMapLink[TKey](cur, c.linkField)
	if u.left != nil {
		walk := u.left
		for {
			v := getMapLink[TKey](walk, c.linkField)
			if v.right == nil {
				break
			}
			walk = v.right
		}
		return walk
	}

	p := u.parent
	for p != nil && getMapLink[TKey](p, c.linkField).left == cur {
		cur = p
		p = getMapLink[TKey](cur, c.linkField).parent
	}
	return p
}

func (c *embeddedMap[TKey, T]) GetPosition(obj *T) int {
	walk := obj
	prev := walk
	position := 0
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if u.left != prev {
			if u.left != nil {
				position += getMapLink[TKey](u.left, c.linkField).position
			}
			position++
		}
		prev = walk
		walk = u.parent
	}
	return position - 1
}

func (c *embeddedMap[TKey, T]) Position(index int) *T {
	walk := c.root
	walkIndex := 0
	if walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if u.left != nil {
			walkIndex = getMapLink[TKey](u.left, c.linkField).position
		}
	}
	for walk != nil {
		u := getMapLink[TKey](walk, c.linkField)
		if index < walkIndex {
			walk = u.left
			walkIndex--
			if u.right != nil {
				walkIndex -= getMapLink[TKey](u.right, c.linkField).position
			}
		} else if walkIndex < index {
			walk = u.right
			walkIndex++
			if u.left != nil {
				walkIndex -= getMapLink[TKey](u.left, c.linkField).position
			}
		} else {
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) Insert(key TKey, obj *T) *T {
	var parent *T
	parentBranch := &c.root
	walk := c.root
	for walk != nil {
		parent = walk
		getMapLink[TKey](parent, c.linkField).position++
		u := getMapLink[TKey](walk, c.linkField)
		if key < u.key {
			parentBranch = &u.left
			walk = *parentBranch
		} else if u.key < key {
			parentBranch = &u.right
			walk = *parentBranch
		}
	}

	*parentBranch = obj
	u := getMapLink[TKey](obj, c.linkField)
	u.parent = parent
	u.left = nil
	u.right = nil
	u.red = true
	u.position = 1
	u.key = key
	c.count++
	c.insertFixup(obj)
	return obj
}

func (c *embeddedMap[TKey, T]) Remove(obj *T) *T {
	u := getMapLink[TKey](obj, c.linkField)
	if u.left != nil && u.right != nil {
		succ := c.Next(obj)
		curParent := u.parent
		curLeft := u.left
		curRight := u.right
		curRed := u.red
		curParentChild := &c.root
		if curParent != nil {
			v := getMapLink[TKey](curParent, c.linkField)
			if v.left == obj {
				curParentChild = &v.left
			} else {
				curParentChild = &v.right
			}
		}

		v := getMapLink[TKey](succ, c.linkField)
		succParent := v.parent
		succLeft := v.left
		succRight := v.right
		succRed := v.red
		w := getMapLink[TKey](succParent, c.linkField)
		succParentChild := &w.right
		if w.left == succ {
			succParentChild = &w.left
		}

		u.position, v.position = v.position, u.position

		u.left, u.right, u.red = succLeft, succRight, succRed
		v.parent = curParent
		v.left, v.right, v.red = curLeft, curRight, curRed
		u.parent = succ
		*curParentChild = succ
		if succRight != nil {
			getMapLink[TKey](succRight, c.linkField).parent = obj
		}

		if succParent == obj {
			u.parent = succ
			v.right = obj
		} else {
			u.parent = succParent
			*succParentChild = obj
			v.right = curRight
			getMapLink[TKey](curRight, c.linkField).parent = succ
		}
	}

	if u.red {
		c.cutNode(obj)
	} else {
		if u.left == nil && u.right == nil {
			c.removeFixup(obj)
			c.cutNode(obj)
		} else {
			child := u.left
			if u.right != nil {
				child = u.right
			}
			x := getMapLink[TKey](child, c.linkField)
			if x.red {
				x.red = false
				c.cutNode(obj)
			} else {
				c.cutNode(obj)
				c.removeFixup(child)
			}
		}
	}

	c.count--
	return obj
}

func (c *embeddedMap[TKey, T]) RemoveFirst() *T {
	head := c.First()
	if head == nil {
		return nil
	}
	return c.Remove(head)
}

func (c *embeddedMap[TKey, T]) RemoveLast() *T {
	tail := c.Last()
	if tail == nil {
		return nil
	}
	return c.Remove(tail)
}

func (c *embeddedMap[TKey, T]) Move(cur *T, newKey TKey) {
	c.Remove(cur)
	c.Insert(newKey, cur)
}

func (c *embeddedMap[TKey, T]) GetKey(obj *T) TKey {
	return getMapLink[TKey](obj, c.linkField).key
}

func (c *embeddedMap[TKey, T]) Count() int {
	return c.count
}

func (c *embeddedMap[TKey, T]) IsEmpty() bool {
	return c.count == 0
}

func (c *embeddedMap[TKey, T]) RemoveAll() {
	c.root = nil
	c.count = 0
}

func (c *embeddedMap[TKey, T]) IsContained(obj *T) bool {
	walk := c.FindFirst(c.GetKey(obj))
	for walk != nil {
		if walk == obj {
			return true
		}
		walk = c.FindNext(walk)
	}
	return false
}

func (c *embeddedMap[TKey, T]) insertFixup(cur *T) {
	u := getMapLink[TKey](cur, c.linkField)
	if u.parent == nil {
		u.red = false
	} else if v := getMapLink[TKey](u.parent, c.linkField); v.red {
		parent := u.parent
		grand := v.parent
		w := getMapLink[TKey](grand, c.linkField)
		uncle := w.right
		if w.left == parent {
			uncle = w.left
		}

		var x *MapLink[TKey, T]
		if uncle != nil {
			x = getMapLink[TKey](uncle, c.linkField)
		}

		if x != nil && x.red {
			v.red = false
			x.red = false
			w.red = true
			c.insertFixup(grand)
		} else {
			if cur == v.right && parent == w.left {
				c.rotateLeft(parent)
				cur = u.left
			} else if cur == v.left && parent == w.right {
				c.rotateRight(parent)
				cur = u.right
			}

			parent = u.parent
			v = getMapLink[TKey](parent, c.linkField)
			grand = v.parent
			w = getMapLink[TKey](grand, c.linkField)
			v.red = false
			w.red = true
			if cur == v.left && parent == w.left {
				c.rotateRight(grand)
			} else {
				c.rotateLeft(grand)
			}
		}
	}
}

func (c *embeddedMap[TKey, T]) rotateLeft(cur *T) {
	u := getMapLink[TKey](cur, c.linkField)
	r := u.right
	p := u.parent
	b := u.left

	u.parent = r
	u.right = b

	v := getMapLink[TKey](r, c.linkField)
	v.parent = p
	v.left = cur

	newC := u.position - v.position
	if b != nil {
		w := getMapLink[TKey](b, c.linkField)
		w.parent = cur
		newC += w.position
	}
	v.position = u.position
	u.position = newC

	if p == nil {
		c.root = r
	} else if x := getMapLink[TKey](p, c.linkField); x.left == cur {
		x.left = r
	} else {
		x.right = r
	}
}

func (c *embeddedMap[TKey, T]) rotateRight(cur *T) {
	u := getMapLink[TKey](cur, c.linkField)
	l := u.left
	p := u.parent
	b := u.right

	u.parent = l
	u.left = b

	v := getMapLink[TKey](l, c.linkField)
	v.parent = p
	v.right = cur

	newC := u.position - v.position
	if b != nil {
		w := getMapLink[TKey](b, c.linkField)
		w.parent = cur
		newC += w.position
	}
	v.position = u.position
	u.position = newC

	if p == nil {
		c.root = l
	} else if x := getMapLink[TKey](p, c.linkField); x.right == cur {
		x.right = l
	} else {
		x.left = l
	}
}

func (c *embeddedMap[TKey, T]) removeFixup(cur *T) {
	u := getMapLink[TKey](cur, c.linkField)
	parent := u.parent
	if parent == nil {
		return
	}

	parentU := getMapLink[TKey](parent, c.linkField)
	sibling := parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	if sibling != nil {
		w := getMapLink[TKey](sibling, c.linkField)
		if w.red {
			parentU.red = true
			w.red = false
			if parentU.left == cur {
				c.rotateLeft(parent)
			} else {
				c.rotateRight(parent)
			}
		}
	}

	parent = u.parent
	parentU = getMapLink[TKey](parent, c.linkField)
	sibling = parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	siblingU := getMapLink[TKey](sibling, c.linkField)
	sibLeft := siblingU.left
	sibRight := siblingU.right

	var sibLeftU *MapLink[TKey, T]
	if sibLeft != nil {
		sibLeftU = getMapLink[TKey](sibLeft, c.linkField)
	}

	var sibRightU *MapLink[TKey, T]
	if sibRight != nil {
		sibRightU = getMapLink[TKey](sibRight, c.linkField)
	}

	if !parentU.red && !siblingU.red && (sibLeftU == nil || !sibLeftU.red) && (sibRightU == nil || !sibRightU.red) {
		siblingU.red = true
		c.removeFixup(parent)
		return
	}

	if parentU.red && !siblingU.red && (sibLeftU == nil || !sibLeftU.red) && (sibRightU == nil || !sibRightU.red) {
		siblingU.red = true
		parentU.red = false
		return
	}

	if cur == parentU.left && !siblingU.red && (sibLeftU != nil && sibLeftU.red) && (sibRightU == nil || !sibRightU.red) {
		siblingU.red = true
		sibLeftU.red = false
		c.rotateRight(sibling)
	} else if cur == parentU.right && !siblingU.red && (sibRightU != nil && sibRightU.red) && (sibLeftU == nil || !sibLeftU.red) {
		siblingU.red = true
		sibRightU.red = false
		c.rotateLeft(sibling)
	}

	parent = u.parent
	parentU = getMapLink[TKey](parent, c.linkField)
	sibling = parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	siblingU = getMapLink[TKey](sibling, c.linkField)
	siblingU.red = parentU.red
	parentU.red = false
	if cur == parentU.left {
		sibRight = siblingU.right
		sibRightU = getMapLink[TKey](sibRight, c.linkField)
		sibRightU.red = false
		c.rotateLeft(parent)
		return
	}

	sibLeft = siblingU.left
	sibLeftU = getMapLink[TKey](sibLeft, c.linkField)
	sibLeftU.red = false
	c.rotateRight(parent)
}

func (c *embeddedMap[TKey, T]) cutNode(cur *T) {
	u := getMapLink[TKey](cur, c.linkField)
	child := u.left
	if u.left == nil {
		child = u.right
	}
	parent := u.parent
	if parent == nil {
		c.root = child
	} else if parentU := getMapLink[TKey](parent, c.linkField); parentU.left == cur {
		parentU.left = child
	} else {
		parentU.right = child
	}

	if child != nil {
		childU := getMapLink[TKey](child, c.linkField)
		childU.parent = parent
	}

	walk := parent
	for walk != nil {
		walkU := getMapLink[TKey](walk, c.linkField)
		walkU.position--
		walk = walkU.parent
	}
}
