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
	// do nothing
}

func (a *staticArray[T]) Size() int {
	return len(a.data)
}

func (a *staticArray[T]) Slice() []T {
	return a.data
}
