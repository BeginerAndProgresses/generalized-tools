package gttype

type LinkList[T any] interface {
	Find(n int) T        // 寻找第N个元素
	Insert(n int, val T) // 在第N处插入元素
	Delete(n int) T      // 删除第N个元素
	Size() int           // 删除某个元素
}

// LinkList 双向链表
type adkLinkList[T any] struct {
	head   *node[T]
	length int
}

type node[T any] struct {
	pri, aft *node[T]
	val      T
}

func NewLinkList[T any]() LinkList[T] {
	nnode := new(node[T])
	nnode.aft = nnode
	nnode.pri = nnode
	return &adkLinkList[T]{head: nnode, length: 0}
}

func (a *adkLinkList[T]) Find(n int) T {
	if n <= 0 && n > a.length {
		return *new(T)
	}
	p := a.head
	for n > 0 {
		p = p.aft
		n--
	}
	return p.val
}

func (a *adkLinkList[T]) Insert(n int, val T) {
	if n < 0 && n > a.length {
		return
	}
	p := a.head
	for n > 0 {
		p = p.aft
		n--
	}
	nnode := &node[T]{val: val}
	p.aft.pri = nnode
	nnode.aft = p.aft
	p.aft = nnode
	nnode.pri = p
	a.length++
	return
}

func (a *adkLinkList[T]) Delete(n int) T {
	if n < 0 && n > a.length {
		return *new(T)
	}
	p := a.head
	for n-1 > 0 {
		p = p.aft
		n--
	}
	val := p.aft.val
	p.aft.aft.pri = p
	p.aft = p.aft.aft
	a.length--
	return val
}

func (a *adkLinkList[T]) Size() int {
	return a.length
}

func (a *adkLinkList[T]) IsEmpty() bool {
	return a.length == 0
}
