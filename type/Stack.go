package gttype

type Stack[T any] interface {
	Pop() T
	Push(val T)
	Top() T
	IsEmpty() bool
	Size() int
}

type sdkStack[T any] struct {
	stack []T
	idx   int
}

func (s *sdkStack[T]) Pop() T {
	if s.IsEmpty() {
		return *new(T)
	}
	top := s.stack[s.idx]
	s.stack = s.stack[:s.idx]
	s.idx--
	return top
}

func (s *sdkStack[T]) Push(val T) {
	s.idx++
	s.stack = append(s.stack, val)
}

func (s *sdkStack[T]) Top() T {
	if s.IsEmpty() {
		return *new(T)
	}
	return s.stack[s.idx]
}

func (s *sdkStack[T]) IsEmpty() bool {
	if s.idx == -1 {
		return true
	}
	return false
}

func (s *sdkStack[T]) Size() int {
	return s.idx + 1
}

// NewStack 初始容量为32，防止频繁扩容影响性能
func NewStack[T any]() Stack[T] {
	return &sdkStack[T]{stack: make([]T, 0, 32), idx: -1}
}
