package embedded

import (
	"constraints"
)

type Map[TKey MapKeyType, T any] interface {
	First() *T
	Last() *T
	Next(cur *T) *T
	Prev(cur *T) *T
	Position(index int) *T
	Count() int

	Remove(obj *T) *T
	RemoveFirst() *T
	RemoveLast() *T
	RemoveAll()

	Insert(key TKey, obj *T) *T

	Move(obj *T, newKey TKey)

	GetKey(obj *T) TKey
	IsEmpty() bool

	IsContained(obj *T) bool

	GetPosition(obj *T) int

	Find(key TKey) *T
	FindFirst(key TKey) *T
	FindNext(cur *T) *T
	FindLowerInclusive(key TKey) *T
	FindUpperInclusive(key TKey) *T
	FindLowerExclusive(key TKey) *T
	FindUpperExclusive(key TKey) *T
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

func (c *embeddedMap[TKey, T]) getLink(obj *T) *MapLink[TKey, T] {
	return getMapLink[TKey](obj, c.linkField)
}

func (c *embeddedMap[TKey, T]) Find(key TKey) *T {
	walk := c.root
	for walk != nil {
		walkU := c.getLink(walk)
		if key < walkU.key {
			walk = walkU.left
		} else if walkU.key < key {
			walk = walkU.right
		} else {
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindFirst(key TKey) *T {
	walk := c.root
	for walk != nil {
		walkU := c.getLink(walk)
		if key < walkU.key {
			walk = walkU.left
		} else if walkU.key < key {
			walk = walkU.right
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
		walkU := c.getLink(walk)
		if key < walkU.key {
			left := walkU.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if walkU.key < key {
			right := walkU.right
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
		walkU := c.getLink(walk)
		if key < walkU.key {
			left := walkU.left
			if left == nil {
				return walk
			}
			walk = left
		} else if walkU.key < key {
			right := walkU.right
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
		walkU := c.getLink(walk)
		if key < walkU.key {
			left := walkU.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if walkU.key < key {
			right := walkU.right
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
		walkU := c.getLink(walk)
		if key < walkU.key {
			left := walkU.left
			if left == nil {
				return walk
			}
			walk = left
		} else if walkU.key < key {
			right := walkU.right
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
		cur = c.getLink(cur).left
	}
	return prev
}

func (c *embeddedMap[TKey, T]) Last() *T {
	var prev *T
	cur := c.root
	for cur != nil {
		prev = cur
		cur = c.getLink(cur).right
	}
	return prev
}

func (c *embeddedMap[TKey, T]) Next(cur *T) *T {
	curU := c.getLink(cur)
	if curU.right != nil {
		walk := curU.right
		for {
			walkU := c.getLink(walk)
			if walkU.left == nil {
				break
			}
			walk = walkU.left
		}
		return walk
	}

	p := curU.parent
	for p != nil && c.getLink(p).right == cur {
		cur = p
		p = c.getLink(cur).parent
	}
	return p
}

func (c *embeddedMap[TKey, T]) Prev(cur *T) *T {
	curU := c.getLink(cur)
	if curU.left != nil {
		walk := curU.left
		for {
			walkU := c.getLink(walk)
			if walkU.right == nil {
				break
			}
			walk = walkU.right
		}
		return walk
	}

	p := curU.parent
	for p != nil && c.getLink(p).left == cur {
		cur = p
		p = c.getLink(cur).parent
	}
	return p
}

func (c *embeddedMap[TKey, T]) GetPosition(obj *T) int {
	walk := obj
	prev := walk
	position := 0
	for walk != nil {
		waklU := c.getLink(walk)
		if waklU.left != prev {
			if waklU.left != nil {
				position += c.getLink(waklU.left).position
			}
			position++
		}
		prev = walk
		walk = waklU.parent
	}
	return position - 1
}

func (c *embeddedMap[TKey, T]) Position(index int) *T {
	walk := c.root
	walkIndex := 0
	if walk != nil {
		walkU := c.getLink(walk)
		if walkU.left != nil {
			walkIndex = c.getLink(walkU.left).position
		}
	}
	for walk != nil {
		walkU := c.getLink(walk)
		if index < walkIndex {
			walk = walkU.left
			walkIndex--
			if walkU.right != nil {
				walkIndex -= c.getLink(walkU.right).position
			}
		} else if walkIndex < index {
			walk = walkU.right
			walkIndex++
			if walkU.left != nil {
				walkIndex -= c.getLink(walkU.left).position
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
		c.getLink(parent).position++
		walkU := c.getLink(walk)
		if key < walkU.key {
			parentBranch = &walkU.left
			walk = *parentBranch
		} else if walkU.key < key {
			parentBranch = &walkU.right
			walk = *parentBranch
		}
	}

	*parentBranch = obj
	objU := c.getLink(obj)
	objU.parent = parent
	objU.left = nil
	objU.right = nil
	objU.red = true
	objU.position = 1
	objU.key = key
	c.count++
	c.insertFixup(obj)
	return obj
}

func (c *embeddedMap[TKey, T]) Remove(obj *T) *T {
	objU := c.getLink(obj)
	if objU.left != nil && objU.right != nil {
		succ := c.Next(obj)
		curParent := objU.parent
		curLeft := objU.left
		curRight := objU.right
		curRed := objU.red
		curParentChild := &c.root
		if curParent != nil {
			curParentU := c.getLink(curParent)
			if curParentU.left == obj {
				curParentChild = &curParentU.left
			} else {
				curParentChild = &curParentU.right
			}
		}

		succU := c.getLink(succ)
		succParent := succU.parent
		succLeft := succU.left
		succRight := succU.right
		succRed := succU.red
		succParentU := c.getLink(succParent)
		succParentChild := &succParentU.right
		if succParentU.left == succ {
			succParentChild = &succParentU.left
		}

		objU.position, succU.position = succU.position, objU.position

		objU.left, objU.right, objU.red = succLeft, succRight, succRed
		succU.parent = curParent
		succU.left, succU.right, succU.red = curLeft, curRight, curRed
		objU.parent = succ
		*curParentChild = succ
		if succRight != nil {
			c.getLink(succRight).parent = obj
		}

		if succParent == obj {
			objU.parent = succ
			succU.right = obj
		} else {
			objU.parent = succParent
			*succParentChild = obj
			succU.right = curRight
			c.getLink(curRight).parent = succ
		}
	}

	if objU.red {
		c.cutNode(obj)
	} else {
		if objU.left == nil && objU.right == nil {
			c.removeFixup(obj)
			c.cutNode(obj)
		} else {
			child := objU.left
			if objU.right != nil {
				child = objU.right
			}
			childU := c.getLink(child)
			if childU.red {
				childU.red = false
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
	return c.getLink(obj).key
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
	curU := c.getLink(cur)
	if curU.parent == nil {
		curU.red = false
	} else if parentU := c.getLink(curU.parent); parentU.red {
		parent := curU.parent
		grand := parentU.parent
		grandU := c.getLink(grand)
		uncle := grandU.right
		if grandU.left == parent {
			uncle = grandU.left
		}

		var uncleU *MapLink[TKey, T]
		if uncle != nil {
			uncleU = c.getLink(uncle)
		}

		if uncleU != nil && uncleU.red {
			parentU.red = false
			uncleU.red = false
			grandU.red = true
			c.insertFixup(grand)
		} else {
			if cur == parentU.right && parent == grandU.left {
				c.rotateLeft(parent)
				cur = curU.left
			} else if cur == parentU.left && parent == grandU.right {
				c.rotateRight(parent)
				cur = curU.right
			}

			parent = curU.parent
			parentU = c.getLink(parent)
			grand = parentU.parent
			grandU = c.getLink(grand)
			parentU.red = false
			grandU.red = true
			if cur == parentU.left && parent == grandU.left {
				c.rotateRight(grand)
			} else {
				c.rotateLeft(grand)
			}
		}
	}
}

func (c *embeddedMap[TKey, T]) rotateLeft(cur *T) {
	curU := c.getLink(cur)
	right := curU.right
	parent := curU.parent
	left := curU.left

	curU.parent = right
	curU.right = left

	rightU := c.getLink(right)
	rightU.parent = parent
	rightU.left = cur

	newC := curU.position - rightU.position
	if left != nil {
		leftU := c.getLink(left)
		leftU.parent = cur
		newC += leftU.position
	}
	rightU.position = curU.position
	curU.position = newC

	if parent == nil {
		c.root = right
	} else if parentU := c.getLink(parent); parentU.left == cur {
		parentU.left = right
	} else {
		parentU.right = right
	}
}

func (c *embeddedMap[TKey, T]) rotateRight(cur *T) {
	curU := c.getLink(cur)
	left := curU.left
	parent := curU.parent
	right := curU.right

	curU.parent = left
	curU.left = right

	leftU := c.getLink(left)
	leftU.parent = parent
	leftU.right = cur

	newC := curU.position - leftU.position
	if right != nil {
		rightU := c.getLink(right)
		rightU.parent = cur
		newC += rightU.position
	}
	leftU.position = curU.position
	curU.position = newC

	if parent == nil {
		c.root = left
	} else if parentU := c.getLink(parent); parentU.right == cur {
		parentU.right = left
	} else {
		parentU.left = left
	}
}

func (c *embeddedMap[TKey, T]) removeFixup(cur *T) {
	curU := c.getLink(cur)
	parent := curU.parent
	if parent == nil {
		return
	}

	parentU := c.getLink(parent)
	sibling := parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	if sibling != nil {
		siblingU := c.getLink(sibling)
		if siblingU.red {
			parentU.red = true
			siblingU.red = false
			if parentU.left == cur {
				c.rotateLeft(parent)
			} else {
				c.rotateRight(parent)
			}
		}
	}

	parent = curU.parent
	parentU = c.getLink(parent)
	sibling = parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	siblingU := c.getLink(sibling)
	sibLeft := siblingU.left
	sibRight := siblingU.right

	var sibLeftU *MapLink[TKey, T]
	if sibLeft != nil {
		sibLeftU = c.getLink(sibLeft)
	}

	var sibRightU *MapLink[TKey, T]
	if sibRight != nil {
		sibRightU = c.getLink(sibRight)
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

	parent = curU.parent
	parentU = c.getLink(parent)
	sibling = parentU.left
	if parentU.left == cur {
		sibling = parentU.right
	}
	siblingU = c.getLink(sibling)
	siblingU.red = parentU.red
	parentU.red = false
	if cur == parentU.left {
		sibRight = siblingU.right
		sibRightU = c.getLink(sibRight)
		sibRightU.red = false
		c.rotateLeft(parent)
		return
	}

	sibLeft = siblingU.left
	sibLeftU = c.getLink(sibLeft)
	sibLeftU.red = false
	c.rotateRight(parent)
}

func (c *embeddedMap[TKey, T]) cutNode(cur *T) {
	curU := c.getLink(cur)
	child := curU.left
	if curU.left == nil {
		child = curU.right
	}
	parent := curU.parent
	if parent == nil {
		c.root = child
	} else if parentU := c.getLink(parent); parentU.left == cur {
		parentU.left = child
	} else {
		parentU.right = child
	}

	if child != nil {
		childU := c.getLink(child)
		childU.parent = parent
	}

	walk := parent
	for walk != nil {
		walkU := c.getLink(walk)
		walkU.position--
		walk = walkU.parent
	}
}
