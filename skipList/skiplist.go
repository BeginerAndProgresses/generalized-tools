package skipList

import (
	"math/rand"
	"time"
)

// DefaultMaxLevel 是所有新创建的skipList的默认级别。
// 可以全局更改。更改它不会影响现有列表。
// 所有的skipList在创建后都可以通过`SetMaxLevel()`方法更新最大级别。
var DefaultMaxLevel = 48

// preallocDefaultMaxLevel 是 Set new element 时在堆栈上分配内存的常量。
const preallocDefaultMaxLevel = 48

// SkipList is the header of a skip list.
type SkipList struct {
	elementHeader

	comparable Comparable
	rand       *rand.Rand

	maxLevel int
	length   int
	back     *Element
}

// New 将创建一个新的 skip list，其中包含 comparable to compare 键。
// 有很多预定义的严格类型键，如 Int、Float64、String 等。
// 我们可以通过实现 Comparable 接口来创建自定义 Comparitation。
// 如果失败，返回 nil。
func New(comparable Comparable) *SkipList {
	if DefaultMaxLevel <= 0 {
		return nil
	}

	source := rand.NewSource(time.Now().UnixNano())
	return &SkipList{
		elementHeader: elementHeader{
			levels: make([]*Element, DefaultMaxLevel),
		},

		comparable: comparable,
		rand:       rand.New(source),

		maxLevel: DefaultMaxLevel,
	}
}

// Init 重置列表,并删除所有元素。
func (list *SkipList) Init() *SkipList {
	list.back = nil
	list.length = 0
	list.levels = make([]*Element, len(list.levels))
	return list
}

// SetRandSource 设置新的 rand 源。
// 默认情况下，Skiplist 使用 math/rand 中定义的全局 rand。
// 默认 rand 在生成任何数字之前获取全局互斥锁。
// 如果 skiplist 受到 caller 的良好保护，则没有必要设置。
func (list *SkipList) SetRandSource(source rand.Source) {
	list.rand = rand.New(source)
}

// Front 返回第一个元素。
func (list *SkipList) Front() (front *Element) {
	return list.levels[0]
}

// Back 返回最后一个元素。
func (list *SkipList) Back() *Element {
	return list.back
}

// Len 返回列表的长度。
func (list *SkipList) Len() int {
	return list.length
}

// Set 设置 键为 key 的元素为 value。
func (list *SkipList) Set(key, value interface{}) (elem *Element) {
	score := list.calcScore(key)

	// 预期列表为空
	if list.length == 0 {
		level := list.randLevel()
		elem = newElement(list, level, score, key, value)

		for i := 0; i < level; i++ {
			list.levels[i] = elem
		}

		list.back = elem
		list.length++
		return
	}

	//找到每一层可能的前置元素
	max := len(list.levels)
	prevHeader := &list.elementHeader

	var maxStaticAllocElemHeaders [preallocDefaultMaxLevel]*elementHeader
	var prevElemHeaders []*elementHeader

	if max <= preallocDefaultMaxLevel {
		prevElemHeaders = maxStaticAllocElemHeaders[:max]
	} else {
		prevElemHeaders = make([]*elementHeader, max)
	}

	for i := max - 1; i >= 0; {
		prevElemHeaders[i] = prevHeader

		for next := prevHeader.levels[i]; next != nil; next = prevHeader.levels[i] {
			if comp := list.compare(score, key, next); comp <= 0 {
				//找到具有相同 key 的 elem。
				//Update value 并返回 elem。
				if comp == 0 {
					elem = next
					elem.Value = value
					return
				}

				break
			}

			prevHeader = &next.elementHeader
			prevElemHeaders[i] = prevHeader
		}

		//如果它们指向与 topLevel 相同的元素，则跳过级别。
		topLevel := prevHeader.levels[i]
		for i--; i >= 0 && prevHeader.levels[i] == topLevel; i-- {
			prevElemHeaders[i] = prevHeader
		}
	}

	// 创建一个新元素
	level := list.randLevel()
	elem = newElement(list, level, score, key, value)

	// 设置一个前置值
	if prev := prevElemHeaders[0]; prev != &list.elementHeader {
		elem.prev = prev.Element()
	}

	// 设置级别
	if prev := prevElemHeaders[level-1]; prev != &list.elementHeader {
		elem.prevTopLevel = prev.Element()
	}

	// 设置级别
	for i := 0; i < level; i++ {
		elem.levels[i] = prevElemHeaders[i].levels[i]
		prevElemHeaders[i].levels[i] = elem
	}

	//找出带有 next 元素的最大级别。
	largestLevel := 0

	for i := level - 1; i >= 0; i-- {
		if elem.levels[i] != nil {
			largestLevel = i + 1
			break
		}
	}

	//调整 next 元素的 prev 和 prevTopLevel。
	if next := elem.levels[0]; next != nil {
		next.prev = elem
	}

	for i := 0; i < largestLevel; {
		next := elem.levels[i]
		nextLevel := next.Level()

		if nextLevel <= level {
			next.prevTopLevel = elem
		}

		i = nextLevel
	}

	//如果 elem 是最后一个元素，则将其设置为 back
	if elem.Next() == nil {
		list.back = elem
	}

	list.length++
	return
}

func (list *SkipList) findNext(start *Element, score float64, key interface{}) (elem *Element) {
	if list.length == 0 {
		return
	}

	if start == nil && list.compare(score, key, list.Front()) <= 0 {
		elem = list.Front()
		return
	}
	if start != nil && list.compare(score, key, start) <= 0 {
		elem = start
		return
	}
	if list.compare(score, key, list.Back()) > 0 {
		return
	}

	var prevHeader *elementHeader
	if start == nil {
		prevHeader = &list.elementHeader
	} else {
		prevHeader = &start.elementHeader
	}
	i := len(prevHeader.levels) - 1

	for i >= 0 {
		for next := prevHeader.levels[i]; next != nil; next = prevHeader.levels[i] {
			if comp := list.compare(score, key, next); comp <= 0 {
				elem = next
				if comp == 0 {
					return
				}

				break
			}

			prevHeader = &next.elementHeader
		}

		topLevel := prevHeader.levels[i]

		for i--; i >= 0 && prevHeader.levels[i] == topLevel; i-- {
		}
	}

	return
}

// FindNext 返回 start 后大于或等于 key 的第一个元素。
// 如果 start 大于或等于 key，则返回 start。
// 如果没有此类元素，则返回 nil。
// 如果 start 为 nil，则从前面查找元素。
func (list *SkipList) FindNext(start *Element, key interface{}) (elem *Element) {
	return list.findNext(start, list.calcScore(key), key)
}

// Find 返回大于或等于 key 的第一个元素。
// 它是 FindNext(nil,key) 的简写。
func (list *SkipList) Find(key interface{}) (elem *Element) {
	return list.FindNext(nil, key)
}

// Get 返回一个带有 key 的元素。
// 如果未找到 key，则返回 nil。
func (list *SkipList) Get(key interface{}) (elem *Element) {
	score := list.calcScore(key)

	firstElem := list.findNext(nil, score, key)
	if firstElem == nil {
		return
	}

	if list.compare(score, key, firstElem) != 0 {
		return
	}

	elem = firstElem
	return
}

// GetValue 返回具有键的元素的值。
func (list *SkipList) GetValue(key interface{}) (val interface{}, ok bool) {
	element := list.Get(key)

	if element == nil {
		return
	}

	val = element.Value
	ok = true
	return
}

// MustGetValue 返回 key 对应的元素
// 如果list 中不存在 则return nil
func (list *SkipList) MustGetValue(key interface{}) interface{} {
	element := list.Get(key)

	if element == nil {
		return nil
	}

	return element.Value
}

// Remove 删除元素。
// 如果找到，则返回已删除的元素指针，如果未找到，则返回 nil。
func (list *SkipList) Remove(key interface{}) (elem *Element) {
	elem = list.Get(key)

	if elem == nil {
		return
	}

	list.RemoveElement(elem)
	return
}

// RemoveFront 删除 front element 节点并返回已删除的元素。
func (list *SkipList) RemoveFront() (front *Element) {
	if list.length == 0 {
		return
	}

	front = list.Front()
	list.RemoveElement(front)
	return
}

// RemoveBack 删除后面的元素节点并返回被删除的元素。
func (list *SkipList) RemoveBack() (back *Element) {
	if list.length == 0 {
		return
	}

	back = list.back
	list.RemoveElement(back)
	return
}

// RemoveElement 从列表中删除 elem。
func (list *SkipList) RemoveElement(elem *Element) {
	if elem == nil || elem.list != list {
		return
	}

	level := elem.Level()

	max := 0
	prevElems := make([]*Element, level)
	prev := elem.prev

	for prev != nil && max < level {
		prevLevel := len(prev.levels)

		for ; max < prevLevel && max < level; max++ {
			prevElems[max] = prev
		}

		for prev = prev.prevTopLevel; prev != nil && prev.Level() == prevLevel; prev = prev.prevTopLevel {
		}
	}

	for i := 0; i < max; i++ {
		prevElems[i].levels[i] = elem.levels[i]
	}

	for i := max; i < level; i++ {
		list.levels[i] = elem.levels[i]
	}

	if next := elem.Next(); next != nil {
		next.prev = elem.prev
	}

	for i := 0; i < level; {
		next := elem.levels[i]

		if next == nil || next.prevTopLevel != elem {
			break
		}

		i = next.Level()
		next.prevTopLevel = prevElems[i-1]
	}

	if list.back == elem {
		list.back = elem.prev
	}

	list.length--
	elem.reset()
}

// MaxLevel 返回当前 Max Level 值。
func (list *SkipList) MaxLevel() int {
	return list.maxLevel
}

// SetMaxLevel 更改跳过列表最大级别。
// 如果 level 不大于 0，则返回 -1。
func (list *SkipList) SetMaxLevel(level int) (old int) {
	if level <= 0 {
		return -1
	}

	list.maxLevel = level
	old = len(list.levels)

	if level == old {
		return
	}

	if old > level {
		for i := old - 1; i >= level; i-- {
			if list.levels[i] != nil {
				level = i
				break
			}
		}

		list.levels = list.levels[:level]
		return
	}

	if level <= cap(list.levels) {
		list.levels = list.levels[:level]
		return
	}

	levels := make([]*Element, level)
	copy(levels, list.levels)
	list.levels = levels
	return
}

func (list *SkipList) randLevel() int {
	estimated := list.maxLevel
	const prob = 1 << 30 // Half of 2^31.
	rand := list.rand
	i := 1

	for ; i < estimated; i++ {
		if rand.Int31() < prob {
			break
		}
	}

	return i
}

// compare 比较两个元素的值并返回 -1、0 和 1。
func (list *SkipList) compare(score float64, key interface{}, rhs *Element) int {
	if score != rhs.score {
		if score > rhs.score {
			return 1
		} else if score < rhs.score {
			return -1
		}

		return 0
	}

	return list.comparable.Compare(key, rhs.key)
}

func (list *SkipList) calcScore(key interface{}) (score float64) {
	score = list.comparable.CalcScore(key)
	return
}
