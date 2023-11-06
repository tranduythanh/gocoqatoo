package stack

type Element[T any] struct {
	value T
	next  *Element[T]
}

type Stack[T any] struct {
	top  *Element[T]
	size int
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Len returns the stack's length
func (s *Stack[T]) Len() int {
	return s.size
}

// Push adds a new element onto the stack
func (s *Stack[T]) Push(value T) {
	s.top = &Element[T]{value, s.top}
	s.size++
}

// Pop removes the top element from the stack and returns its value
// If the stack is empty, it returns the zero value for the type T
func (s *Stack[T]) Pop() (value T) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	var zero T // zero value for type T
	return zero
}
