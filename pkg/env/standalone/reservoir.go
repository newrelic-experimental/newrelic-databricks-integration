package standalone

import (
	"math/rand"
)

type ReservoirBuffer[T any] struct {
	buff     []T
	size     int
	capacity int
	index    int
}

func MakeReservoirBuffer[T any](capacity int) Buffer[T] {
	return &ReservoirBuffer[T]{
		buff:     make([]T, capacity),
		size:     0,
		capacity: capacity,
		index:    0,
	}
}

// Implement reservoir sampling Algorithm R
func (b *ReservoirBuffer[T]) Put(val T) bool {
	accepted := false
	if b.index < b.capacity {
		// Fill buffer
		b.buff[b.index] = val
		b.size += 1
		accepted = true
	} else {
		// Buffer is full, start sampling
		j := rand.Intn(b.index)
		if j < b.capacity {
			b.buff[j] = val
			accepted = true
		} else {
			accepted = false
		}
	}
	b.index += 1
	return accepted
}

func (b *ReservoirBuffer[T]) Capacity() int {
	return b.capacity
}

func (b *ReservoirBuffer[T]) Size() int {
	return b.size
}

func (b *ReservoirBuffer[T]) Clear() *[]T {
	ret := b.buff
	b.buff = make([]T, b.capacity)
	b.size = 0
	b.index = 0
	return &ret
}
