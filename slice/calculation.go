package slice

// Intersect 求交集
func Intersect[T comparable](vars1 []T, other ...[]T) []T {
	emCap := len(vars1) / 2
	m := make(map[T]struct{}, emCap)
	res := make([]T, 0, emCap)
	for _, t := range vars1 {
		m[t] = struct{}{}
	}
	for _, ts := range other {
		for _, t := range ts {
			if _, ok := m[t]; ok {
				res = append(res, t)
			}
		}
	}
	return res
}

// Union 求并集
func Union[T comparable](vars ...[]T) []T {
	emCap := len(vars) * 16
	m := make(map[T]struct{}, emCap)
	res := make([]T, 0, emCap)
	for _, ts := range vars {
		for _, t := range ts {
			m[t] = struct{}{}
		}
	}
	for t, _ := range m {
		res = append(res, t)
	}
	return res
}

// Diff 求vars1对vars2的差集
func Diff[T comparable](vars1, vars2 []T) []T {
	emCap := len(vars1) * 4 / 5
	m := make(map[T]struct{}, emCap)
	res := make([]T, 0, emCap)
	for _, t := range vars1 {
		m[t] = struct{}{}
	}
	for _, t := range vars2 {
		if _, ok := m[t]; ok {
			delete(m, t)
		}
	}
	for t, _ := range m {
		res = append(res, t)
	}
	return res
}
