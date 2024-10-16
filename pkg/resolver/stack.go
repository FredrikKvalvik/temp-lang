package resolver

// The are no safetychecks, user must make sure
// that to do proper checks, or else the program might panic
type Stack[T any] []T

func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}
func (s *Stack[T]) Pop() T {
	last := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return last
}

// will panic if peek is used when stack is empty
func (s *Stack[T]) Peek() T {
	return (*s)[len(*s)-1]
}

func (s *Stack[T]) Get(idx int) T {
	return (*s)[idx]
}

func (s *Stack[T]) Size() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
