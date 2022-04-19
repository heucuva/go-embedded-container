package embedded

import (
	"golang.org/x/exp/constraints"
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
	linkField uintptr
	count     int
}

func (c *embeddedMap[TKey, T]) getLink(obj *T) *MapLink[TKey, T] {
	return getMapLink[TKey](obj, c.linkField)
}

func (c *embeddedMap[TKey, T]) Find(key TKey) *T {
	walk := c.root
	for walk != nil {
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			walk = walkLink.left
		} else if walkLink.key < key {
			walk = walkLink.right
		} else {
			return walk
		}
	}
	return nil
}

func (c *embeddedMap[TKey, T]) FindFirst(key TKey) *T {
	walk := c.root
	for walk != nil {
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			walk = walkLink.left
		} else if walkLink.key < key {
			walk = walkLink.right
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
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			left := walkLink.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if walkLink.key < key {
			right := walkLink.right
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
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			left := walkLink.left
			if left == nil {
				return walk
			}
			walk = left
		} else if walkLink.key < key {
			right := walkLink.right
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
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			left := walkLink.left
			if left == nil {
				return c.Prev(walk)
			}
			walk = left
		} else if walkLink.key < key {
			right := walkLink.right
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
		walkLink := c.getLink(walk)
		if key < walkLink.key {
			left := walkLink.left
			if left == nil {
				return walk
			}
			walk = left
		} else if walkLink.key < key {
			right := walkLink.right
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
	curLink := c.getLink(cur)
	if curLink == nil {
		return nil
	}
	if curLink.right != nil {
		walk := curLink.right
		for {
			walkLink := c.getLink(walk)
			if walkLink.left == nil {
				break
			}
			walk = walkLink.left
		}
		return walk
	}

	curParent := curLink.parent
	for curParent != nil && c.getLink(curParent).right == cur {
		cur = curParent
		curParent = c.getLink(cur).parent
	}
	return curParent
}

func (c *embeddedMap[TKey, T]) Prev(cur *T) *T {
	curLink := c.getLink(cur)
	if curLink == nil {
		return nil
	}
	if curLink.left != nil {
		walk := curLink.left
		for {
			walkLink := c.getLink(walk)
			if walkLink.right == nil {
				break
			}
			walk = walkLink.right
		}
		return walk
	}

	curParent := curLink.parent
	for curParent != nil && c.getLink(curParent).left == cur {
		cur = curParent
		curParent = c.getLink(cur).parent
	}
	return curParent
}

func (c *embeddedMap[TKey, T]) GetPosition(obj *T) int {
	walk := obj
	prev := walk
	position := 0
	for walk != nil {
		walkLink := c.getLink(walk)
		if walkLink.left != prev {
			if walkLink.left != nil {
				position += c.getLink(walkLink.left).position
			}
			position++
		}
		prev = walk
		walk = walkLink.parent
	}
	return position - 1
}

func (c *embeddedMap[TKey, T]) Position(index int) *T {
	walk := c.root
	walkIndex := 0
	if walk != nil {
		walkLink := c.getLink(walk)
		if walkLink.left != nil {
			walkIndex = c.getLink(walkLink.left).position
		}
	}
	for walk != nil {
		walkLink := c.getLink(walk)
		if index < walkIndex {
			walk = walkLink.left
			walkIndex--
			if walkLink.right != nil {
				walkIndex -= c.getLink(walkLink.right).position
			}
		} else if walkIndex < index {
			walk = walkLink.right
			walkIndex++
			if walkLink.left != nil {
				walkIndex -= c.getLink(walkLink.left).position
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
	walk := *parentBranch
	for walk != nil {
		parent = walk
		walkLink := c.getLink(walk)
		walkLink.position++
		if key < walkLink.key {
			parentBranch = &walkLink.left
			walk = *parentBranch
		} else {
			parentBranch = &walkLink.right
			walk = *parentBranch
		}
	}

	*parentBranch = obj
	objLink := c.getLink(obj)
	objLink.parent = parent
	objLink.left = nil
	objLink.right = nil
	objLink.red = true
	objLink.position = 1
	objLink.key = key
	c.count++
	c.insertFixup(obj)
	return obj
}

func (c *embeddedMap[TKey, T]) Remove(obj *T) *T {
	objLink := c.getLink(obj)
	if objLink.left != nil && objLink.right != nil {
		succ := c.Next(obj)
		curParent := objLink.parent
		curLeft := objLink.left
		curRight := objLink.right
		curRed := objLink.red
		curParentChild := &c.root
		if curParent != nil {
			curParentLink := c.getLink(curParent)
			if curParentLink.left == obj {
				curParentChild = &curParentLink.left
			} else {
				curParentChild = &curParentLink.right
			}
		}

		succLink := c.getLink(succ)
		succParent := succLink.parent
		succLeft := succLink.left
		succRight := succLink.right
		succRed := succLink.red
		succParentLink := c.getLink(succParent)
		succParentChild := &succParentLink.right
		if succParentLink.left == succ {
			succParentChild = &succParentLink.left
		}

		objLink.position, succLink.position = succLink.position, objLink.position

		objLink.left, objLink.right, objLink.red = succLeft, succRight, succRed
		succLink.parent = curParent
		succLink.left, succLink.right, succLink.red = curLeft, curRight, curRed
		objLink.parent = succ
		*curParentChild = succ
		if succRight != nil {
			c.getLink(succRight).parent = obj
		}

		if succParent == obj {
			objLink.parent = succ
			succLink.right = obj
		} else {
			objLink.parent = succParent
			*succParentChild = obj
			succLink.right = curRight
			c.getLink(curRight).parent = succ
		}
	}

	if objLink.red {
		c.cutNode(obj)
	} else {
		if objLink.left == nil && objLink.right == nil {
			c.removeFixup(obj)
			c.cutNode(obj)
		} else {
			child := objLink.left
			if objLink.right != nil {
				child = objLink.right
			}
			childLink := c.getLink(child)
			if childLink.red {
				childLink.red = false
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
	curLink := c.getLink(cur)
	if curLink.parent == nil {
		curLink.red = false
	} else if parentLink := c.getLink(curLink.parent); parentLink.red {
		parent := curLink.parent
		grand := parentLink.parent
		grandLink := c.getLink(grand)
		var uncle *T
		if grandLink != nil {
			uncle = grandLink.right
			if grandLink.left == parent {
				uncle = grandLink.left
			}
		}

		var uncleLink *MapLink[TKey, T]
		if uncle != nil {
			uncleLink = c.getLink(uncle)
		}

		if uncleLink != nil && uncleLink.red {
			parentLink.red = false
			uncleLink.red = false
			grandLink.red = true
			c.insertFixup(grand)
		} else if grandLink != nil {
			if cur == parentLink.right && parent == grandLink.left {
				c.rotateLeft(parent)
				cur = curLink.left
			} else if cur == parentLink.left && parent == grandLink.right {
				c.rotateRight(parent)
				cur = curLink.right
			}

			parent = curLink.parent
			parentLink = c.getLink(parent)
			grand = parentLink.parent
			grandLink = c.getLink(grand)
			parentLink.red = false
			grandLink.red = true
			if cur == parentLink.left && parent == grandLink.left {
				c.rotateRight(grand)
			} else {
				c.rotateLeft(grand)
			}
		}
	}
}

func (c *embeddedMap[TKey, T]) rotateLeft(cur *T) {
	curLink := c.getLink(cur)
	right := curLink.right
	parent := curLink.parent
	left := curLink.left

	curLink.parent = right
	curLink.right = left

	rightLink := c.getLink(right)
	var rightPos int
	if rightLink != nil {
		rightLink.parent = parent
		rightLink.left = cur
		rightPos = rightLink.position
	}

	newC := curLink.position - rightPos
	if left != nil {
		leftLink := c.getLink(left)
		leftLink.parent = cur
		newC += leftLink.position
	}
	if rightLink != nil {
		rightLink.position = curLink.position
	}
	curLink.position = newC

	if parent == nil {
		c.root = right
	} else if parentLink := c.getLink(parent); parentLink.left == cur {
		parentLink.left = right
	} else {
		parentLink.right = right
	}
}

func (c *embeddedMap[TKey, T]) rotateRight(cur *T) {
	curLink := c.getLink(cur)
	left := curLink.left
	parent := curLink.parent
	right := curLink.right

	curLink.parent = left
	curLink.left = right

	leftLink := c.getLink(left)
	leftLink.parent = parent
	leftLink.right = cur

	newC := curLink.position - leftLink.position
	if right != nil {
		rightLink := c.getLink(right)
		rightLink.parent = cur
		newC += rightLink.position
	}
	leftLink.position = curLink.position
	curLink.position = newC

	if parent == nil {
		c.root = left
	} else if parentLink := c.getLink(parent); parentLink.right == cur {
		parentLink.right = left
	} else {
		parentLink.left = left
	}
}

func (c *embeddedMap[TKey, T]) removeFixup(cur *T) {
	curLink := c.getLink(cur)
	parent := curLink.parent
	if parent == nil {
		return
	}

	parentLink := c.getLink(parent)
	sibling := parentLink.left
	if parentLink.left == cur {
		sibling = parentLink.right
	}
	if sibling != nil {
		siblingLink := c.getLink(sibling)
		if siblingLink.red {
			parentLink.red = true
			siblingLink.red = false
			if parentLink.left == cur {
				c.rotateLeft(parent)
			} else {
				c.rotateRight(parent)
			}
		}
	}

	parent = curLink.parent
	parentLink = c.getLink(parent)
	sibling = parentLink.left
	if parentLink.left == cur {
		sibling = parentLink.right
	}
	siblingLink := c.getLink(sibling)
	sibLeft := siblingLink.left
	sibRight := siblingLink.right

	var sibLeftLink *MapLink[TKey, T]
	if sibLeft != nil {
		sibLeftLink = c.getLink(sibLeft)
	}

	var sibRightLink *MapLink[TKey, T]
	if sibRight != nil {
		sibRightLink = c.getLink(sibRight)
	}

	if !parentLink.red && !siblingLink.red && (sibLeftLink == nil || !sibLeftLink.red) && (sibRightLink == nil || !sibRightLink.red) {
		siblingLink.red = true
		c.removeFixup(parent)
		return
	}

	if parentLink.red && !siblingLink.red && (sibLeftLink == nil || !sibLeftLink.red) && (sibRightLink == nil || !sibRightLink.red) {
		siblingLink.red = true
		parentLink.red = false
		return
	}

	if cur == parentLink.left && !siblingLink.red && (sibLeftLink != nil && sibLeftLink.red) && (sibRightLink == nil || !sibRightLink.red) {
		siblingLink.red = true
		sibLeftLink.red = false
		c.rotateRight(sibling)
	} else if cur == parentLink.right && !siblingLink.red && (sibRightLink != nil && sibRightLink.red) && (sibLeftLink == nil || !sibLeftLink.red) {
		siblingLink.red = true
		sibRightLink.red = false
		c.rotateLeft(sibling)
	}

	parent = curLink.parent
	parentLink = c.getLink(parent)
	sibling = parentLink.left
	if parentLink.left == cur {
		sibling = parentLink.right
	}
	siblingLink = c.getLink(sibling)
	siblingLink.red = parentLink.red
	parentLink.red = false
	if cur == parentLink.left {
		sibRight = siblingLink.right
		sibRightLink = c.getLink(sibRight)
		sibRightLink.red = false
		c.rotateLeft(parent)
		return
	}

	sibLeft = siblingLink.left
	sibLeftLink = c.getLink(sibLeft)
	sibLeftLink.red = false
	c.rotateRight(parent)
}

func (c *embeddedMap[TKey, T]) cutNode(cur *T) {
	curLink := c.getLink(cur)
	child := curLink.left
	if curLink.left == nil {
		child = curLink.right
	}
	parent := curLink.parent
	if parent == nil {
		c.root = child
	} else if parentLink := c.getLink(parent); parentLink.left == cur {
		parentLink.left = child
	} else {
		parentLink.right = child
	}

	if child != nil {
		childLink := c.getLink(child)
		childLink.parent = parent
	}

	walk := parent
	for walk != nil {
		walkLink := c.getLink(walk)
		walkLink.position--
		walk = walkLink.parent
	}
}
