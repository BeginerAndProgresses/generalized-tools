package slice

import (
	"fmt"
	gterr "generalized-tools/error"
)

// Delete 删除idx指向vars切片中的元素
// 并返回删除后的切片
// 缩容情况：
// cap大于512且cap/len大于等于4，cap大于256且cap/len大于等于2
func Delete[T any](idx int, vars []T) ([]T, T, error) {
	var res T
	if idx < 0 || idx > len(vars)-1 {
		return vars, res, fmt.Errorf(gterr.DeleteIndexError)
	}
	res = vars[idx]
	for i := idx; i < len(vars)-1; i++ {
		vars[i] = vars[i+1]
	}
	nVars := shrink(vars, aShrinkageFactor)
	return nVars, res, nil
}

type shrinkageFactor func(l, c int) (newCap int, needShrink bool)

var aShrinkageFactor shrinkageFactor = func(l, c int) (newCap int, needShrink bool) {
	if c <= 32 {
		return c, false
	}
	if c > 512 && c/l >= 4 {
		return c / 2, true
	}
	if c > 256 && c/l >= 2 {
		factor := 0.625
		return int(float64(c) * factor), true
	}
	return c, false
}

func shrink[T any](vars []T, sf shrinkageFactor) []T {
	length := len(vars)
	nc, ns := sf(length, cap(vars))
	if !ns {
		return vars
	}
	nVars := make([]T, 0, nc)
	for i := 0; i < length; i++ {
		nVars = append(nVars, vars[i])
	}
	return nVars
}

// Insert 在指定位置插入一个元素
func Insert[T any](idx int, val T, vars []T) ([]T, error) {
	if idx < 0 || idx > len(vars) {
		return vars, fmt.Errorf(gterr.InsertIndexError)
	}
	res := append(vars, val)
	for i := len(vars); i > idx; i-- {
		res[i] = res[i-1]
	}
	res[idx] = val
	return res, nil
}

// Filter 过滤可比较类型切片中的元素
func Filter[T comparable](vars []T, elements ...T) []T {
	m := make(map[T]struct{}, len(vars)/2)
	res := make([]T, 0, cap(vars))
	for _, option := range elements {
		m[option] = struct{}{}
	}
	for _, t := range vars {
		if _, ok := m[t]; !ok {
			res = append(res, t)
		}
	}
	return shrink(res, aShrinkageFactor)
}

// Find 查找vars中val下标，并返回所有下标
func Find[T comparable](vars []T, val T) []int {
	res := make([]int, 0, len(vars)/2)
	for i, t := range vars {
		if t == val {
			res = append(res, i)
		}
	}
	return res
}

// FindFirst 获取vars中第一个val下标
// 返回-1没有找到
func FindFirst[T comparable](vars []T, val T) int {
	for i := 0; i < len(vars); i++ {
		if vars[i] == val {
			return i
		}
	}
	return -1
}

// FindLast 获取vars中最后一个val下标
// 返回-1没有找到
func FindLast[T comparable](vars []T, val T) int {
	for i := len(vars) - 1; i >= 0; i-- {
		if vars[i] == val {
			return i
		}
	}
	return -1
}
