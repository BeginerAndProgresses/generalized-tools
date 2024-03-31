package slice

import "math"

// Intersect 求交集
func Intersect[T comparable](vars1, vars2 []T) []T {
	minLen := int(math.Min(float64(len(vars1)), float64(len(vars2))))
	m := make(map[T]struct{}, minLen/2)
	res := make([]T, 0, minLen/2)
	for _, t := range vars1 {
		m[t] = struct{}{}
	}
	for _, t := range vars2 {
		if _, ok := m[t]; ok {
			res = append(res, t)
		}
	}
	return res
}
