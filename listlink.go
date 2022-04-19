package embedded

import (
	"unsafe"
)

// ListLink is a link to the list container
type ListLink[M any] struct {
	prev *M
	next *M
}

func (l *ListLink[M]) remove(linkFieldOfs uintptr, head **M, tail **M) bool {
	if !l.isContained(linkFieldOfs, *head) {
		return false
	}
	if l.prev == nil {
		*head = l.next
	} else {
		getListLink(l.prev, linkFieldOfs).next = l.next
	}
	if l.next == nil {
		*tail = l.prev
	} else {
		getListLink(l.next, linkFieldOfs).prev = l.prev
	}

	l.next = nil
	l.prev = nil
	return true
}

func (l *ListLink[M]) isContained(linkFieldOfs uintptr, head *M) bool {
	return l.prev != nil || head == l.getItem(linkFieldOfs)
}

func (l *ListLink[M]) getItem(linkFieldOfs uintptr) *M {
	u := unsafe.Add(unsafe.Pointer(l), (^linkFieldOfs)+1)
	m := (*M)(u)
	return m
}

func getListLink[T any](obj *T, linkFieldOfs uintptr) *ListLink[T] {
	if obj == nil {
		return nil
	}
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*ListLink[T])(u)
}
