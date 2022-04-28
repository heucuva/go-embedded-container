package array

type staticArray[T any] struct {
	data []T
}

func NewStaticArray[T any](arraySize int) Array[T] {
	return &staticArray[T]{
		data: make([]T, arraySize),
	}
}

func (a *staticArray[T]) IsStatic() bool {
	return true
}

func (a *staticArray[T]) Reserve(count int) {
	panic("cannot reserve a count with a static table size")
}

func (a *staticArray[T]) Size() int {
	return len(a.data)
}

func (a *staticArray[T]) Slice() []T {
	return a.data
}

func (a *staticArray[T]) At(idx int) T {
	return a.data[idx]
}

func (a *staticArray[T]) Ptr(idx int) *T {
	return &a.data[idx]
}

func (a *staticArray[T]) Set(idx int, value T) {
	a.data[idx] = value
}
