package vutils

// Stack represents a stack data structure
type Stack[T any] struct {
	elements []T
}

// Push adds an element to the stack
func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

// Pop removes and returns the top element of the stack
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element, true
}
