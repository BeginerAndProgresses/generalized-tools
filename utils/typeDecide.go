package utils

import (
	"reflect"
)

// IsComparable 判断传入的类型是否可比较
func IsComparable(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.IsValid() && rv.Type().Comparable()
}
