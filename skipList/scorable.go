package skipList

import "reflect"

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/9/7 22:00
 */

type Scorable interface {
	Score() float64
}

// CalcScore 计算key的分数
// 必须遵循以下规则
// 1 不同的分数值必须不相同
// 2 如果返回值为正的，k1.Score() <= k2.Score()
// 3 如果返回值为负的，k1.Score() >= k2.Score()
// 4 如果返回值为0，k1.Score() == k2.Score()
func CalcScore(key any) float64 {
	if scorable, ok := key.(Scorable); ok {
		return scorable.Score()
	}
	return calcScore(reflect.ValueOf(key))
}
