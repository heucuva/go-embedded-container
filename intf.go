package embedded

type TableInterface interface {
	Count() int

	RemoveAll()

	GetTableSize() int
	GetTableUsed() int
	Reserve(count int)
	IsEmpty() bool

	IsStatic() bool
}

type ListInterface[TKey any, T any] interface {
	First() *T
	Last() *T
	Next(cur *T) *T
	Prev(cur *T) *T
	Position(index int) *T

	Remove(obj *T) *T
	RemoveFirst() *T
	RemoveLast() *T
	RemoveAll()

	InsertFirst(key TKey, obj *T) *T
	InsertLast(key TKey, obj *T) *T
	InsertAfter(key TKey, prev, obj *T) *T
	InsertBefore(key TKey, after, obj *T) *T

	Move(obj *T, newKey TKey)
	MoveFirst(obj *T)
	MoveLast(obj *T)
	MoveAfter(dest, obj *T)
	MoveBefore(dest, obj *T)

	FindFirst(key TKey) *T
	FindNext(prevResult *T) *T
	GetKey(obj *T) TKey

	IsContained(obj *T) bool
}
