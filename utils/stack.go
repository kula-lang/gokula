package utils

type Stack[T any] []T

func NewStack[T any]() Stack[T] {
	st := make([]T, 0)
	return st
}

func (s *Stack[T]) Clear() {
	(*s) = (*s)[:0]
}

func (s *Stack[T]) Size() int {
	return len(*s)
}

func (s *Stack[T]) Empty() bool {
	return s.Size() == 0
}

func (s *Stack[T]) Push(v T) {
	(*s) = append((*s), v)
}

func (s *Stack[T]) Pop() (t T) {
	if len(*s) == 0 {
		return
	}
	t = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t
}

func (s *Stack[T]) Peek() (t T) {
	if len(*s) == 0 {
		return
	}
	return (*s)[len(*s)-1]
}
