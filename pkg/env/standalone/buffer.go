package standalone

type Buffer[T any] interface {
	Put(T) bool
	Capacity() int
	Size() int
	Clear() *[]T
}
