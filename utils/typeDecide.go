package utils

import "reflect"

func IsComparable(v any) bool {
	rv := reflect.ValueOf(v)

	return rv.IsValid() && rv.Type().Comparable()
}
