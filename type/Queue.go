package gttype

type Queue[T any] interface {
	Add(val T)
	Element() T
	Remove() T
	IsEmpty() bool
	Size() int
}

type sdkQueue[T any] struct {
	queue []T
}

func (s *sdkQueue[T]) Add(val T) {
	s.queue = append(s.queue, val)
}

func (s *sdkQueue[T]) Element() T {
	if s.IsEmpty() {
		return *new(T)
	}
	return s.queue[0]
}

func (s *sdkQueue[T]) Remove() T {
	if s.IsEmpty() {
		return *new(T)
	} else {
		val := s.queue[0]
		s.queue = s.queue[1:]
		return val
	}
}

func (s *sdkQueue[T]) IsEmpty() bool {
	if len(s.queue) == 0 {
		return true
	}
	return false
}

func (s *sdkQueue[T]) Size() int {
	return len(s.queue)
}

func NewQueue[T any]() Queue[T] {
	return &sdkQueue[T]{queue: make([]T, 0, 32)}
}
