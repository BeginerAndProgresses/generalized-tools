package gttype

type HashSet[T any] interface {
	Add(T) bool
	Size() int
	Remove(T) bool
	IsEmpty() bool
	Clear()
	GetData() []T
}

type adkHashSet[T any] struct {
	Data map[any]struct{}
}

func (a *adkHashSet[T]) GetData() []T {
	var data []T
	for k := range a.Data {
		data = append(data, k)
	}
	return data
}

func (a *adkHashSet[T]) Add(val T) bool {
	if _, ok := a.Data[val]; ok {
		return false
	}
	a.Data[val] = struct{}{}
	return true
}

func (a *adkHashSet[T]) Size() int {
	return len(a.Data)
}

func (a *adkHashSet[T]) Remove(val T) bool {
	delete(a.Data, val)
	return true
}

func (a *adkHashSet[T]) IsEmpty() bool {
	//TODO implement me
	if a.Size() == 0 {
		return true
	}
	return false
}

func (a *adkHashSet[T]) Clear() {
	a.Data = make(map[any]struct{})
}

func NewHashSet[T any]() HashSet[T] {
	return &adkHashSet[T]{
		Data: make(map[any]struct{}),
	}
}
