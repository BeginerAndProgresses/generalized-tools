package gttype

type HashSet[T any] interface {
	Add(T) bool
	Size() int
	Remove(T) bool
	IsEmpty() bool
	Clear()
}

type adkHashSet[T any] struct {
}

func (a *adkHashSet[T]) Add(val T) bool {
	//TODO implement me
	panic("implement me")
}

func (a *adkHashSet[T]) Size() int {
	//TODO implement me
	panic("implement me")
}

func (a *adkHashSet[T]) Remove(val T) bool {
	//TODO implement me
	panic("implement me")
}

func (a *adkHashSet[T]) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (a *adkHashSet[T]) Clear() {
	//TODO implement me
	panic("implement me")
}

func NewHashSet[T any]() HashSet[T] {
	return &adkHashSet[T]{}
}
