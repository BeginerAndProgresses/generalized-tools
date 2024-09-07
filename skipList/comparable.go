package skipList

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/9/2 8:50
 */

// Comparable 比较接口
type Comparable interface {
	// Compare 比较两个对象，返回-1，0，1
	Compare(a, b interface{}) int
	// CalcScore 计算一个对象的分数
	CalcScore(key interface{}) float64
}

type ComparableFunc func(a, b interface{}) int

// Compare 实现Comparable
func (c ComparableFunc) Compare(a, b interface{}) int {
	return c(a, b)
}

// CalcScore 不知道c中func的实现，所以返回0
func (c ComparableFunc) CalcScore(key interface{}) float64 {
	return 0
}
