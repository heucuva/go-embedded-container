package array

type dynamicArray[T any] struct {
	data     []T
	minSize  int
	onResize func(dest, src []T)
}

func NewDynamicArray[T any](minSize int, onResize func(dest, src []T)) Array[T] {
	a := &dynamicArray[T]{
		minSize:  minSize,
		onResize: onResize,
	}
	a.Reserve(minSize)
	return a
}

func (a *dynamicArray[T]) IsStatic() bool {
	return false
}

func (a *dynamicArray[T]) Reserve(count int) {
	count += count >> 2
	if count > len(a.data) {
		a.resize(count)
	}
}

func (a *dynamicArray[T]) Size() int {
	return len(a.data)
}

func (a *dynamicArray[T]) Slice() []T {
	return a.data
}

func (a *dynamicArray[T]) At(idx int) T {
	return a.data[idx]
}

func (a *dynamicArray[T]) Ptr(idx int) *T {
	return &a.data[idx]
}

func (a *dynamicArray[T]) Set(idx int, value T) {
	a.data[idx] = value
}

func (a *dynamicArray[T]) resize(count int) {
	dynamicTableOld := a.data

	dynamicSize := int(nextPowerOf2(uint(count)))
	if dynamicSize < a.minSize {
		dynamicSize = a.minSize
	}
	a.data = make([]T, dynamicSize)
	if a.onResize != nil {
		a.onResize(a.data, dynamicTableOld)
	}
}
