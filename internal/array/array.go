package array

type Array[T any] interface {
	IsStatic() bool
	Reserve(count int)
	Size() int
	Slice() []T
}
