package array

type Array[T any] interface {
	IsStatic() bool
	Reserve(count int)
	Size() int
	Slice() []T
	At(idx int) T
	Ptr(idx int) *T
	Set(idx int, value T)
}
