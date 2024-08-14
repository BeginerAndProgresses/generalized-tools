package slice

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/8/14 21:48
 */

// Convert 将切片中的元素进行转换
// 转换规则为fn
func Convert[T any](vars []T, fn func(T) T) []T {
	var newVars = make([]T, len(vars))
	for i, t := range vars {
		newVars[i] = fn(t)
	}
	return newVars
}

// ToPrt 将切片中的元素转换为指针
func ToPrt[T any](vars []T) []*T {
	var newVars = make([]*T, len(vars))
	for i, t := range vars {
		newVars[i] = &t
	}
	return newVars
}
