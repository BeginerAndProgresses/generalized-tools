package skipList

import "unsafe"

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/9/3 22:02
 */

// 跳表的节点头
// 只能
type elementHeader struct {
	levels []*Element // Next element at all levels.
}

// Element list中的节点
type Element struct {
	elementHeader

	Value interface{}
	key   interface{}
	score float64

	prev         *Element  // Points to previous adjacent elem.
	prevTopLevel *Element  // Points to previous element which points to this element's top most level.
	list         *SkipList // The list contains this elem.
}

func (header *elementHeader) Element() *Element {
	return (*Element)(unsafe.Pointer(header))
}

func newElement(list *SkipList, level int, score float64, key, value interface{}) *Element {
	return &Element{
		elementHeader: elementHeader{
			levels: make([]*Element, level),
		},
		Value: value,
		key:   key,
		score: score,
		list:  list,
	}
}

// Next 返回下一个 elem.
func (elem *Element) Next() *Element {
	if len(elem.levels) == 0 {
		return nil
	}

	return elem.levels[0]
}

// Prev 返回前一个 elem.
func (elem *Element) Prev() *Element {
	return elem.prev
}

// NextLevel 返回特定级别的下一个元素。
// 如果 level 无效，则返回 nil
func (elem *Element) NextLevel(level int) *Element {
	if level < 0 || level >= len(elem.levels) {
		return nil
	}

	return elem.levels[level]
}

// PrevLevel 返回指向特定级别的此元素的上一个元素。
// 如果 level 无效，则返回 nil。
func (elem *Element) PrevLevel(level int) *Element {
	if level < 0 || level >= len(elem.levels) {
		return nil
	}

	if level == 0 {
		return elem.prev
	}

	if level == len(elem.levels)-1 {
		return elem.prevTopLevel
	}

	prev := elem.prev

	for prev != nil {
		if level < len(prev.levels) {
			return prev
		}

		prev = prev.prevTopLevel
	}

	return prev
}

// Key 返回 elem 的 key。
func (elem *Element) Key() interface{} {
	return elem.key
}

// Score 返回此元素的分数。
// Skip list 使所有元素按分数从小到大排序。
func (elem *Element) Score() float64 {
	return elem.score
}

// Level 返回此 elem 的级别。
func (elem *Element) Level() int {
	return len(elem.levels)
}

func (elem *Element) reset() {
	elem.list = nil
	elem.prev = nil
	elem.prevTopLevel = nil
	elem.levels = nil
}
