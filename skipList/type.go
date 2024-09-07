package skipList

import (
	"bytes"
	"fmt"
	"reflect"
)

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/9/2 8:49
 */

// keyType 自定义跳表key值
type keyType int

var typeOfBytes = reflect.TypeOf([]byte(nil))

var _ Comparable = keyType(0)

const (
	byteType    = keyType(reflect.Uint8)
	runeType    = keyType(reflect.Int32)
	intType     = keyType(reflect.Int)
	int8Type    = keyType(reflect.Int8)
	int16Type   = keyType(reflect.Int16)
	int32Type   = keyType(reflect.Int32)
	int64Type   = keyType(reflect.Int64)
	uintType    = keyType(reflect.Uint)
	uint8Type   = keyType(reflect.Uint8)
	uint16Type  = keyType(reflect.Uint16)
	uint32Type  = keyType(reflect.Uint32)
	uint64Type  = keyType(reflect.Uint64)
	uintptrType = keyType(reflect.Uintptr)
	float32Type = keyType(reflect.Float32)
	float64Type = keyType(reflect.Float64)
	stringType  = keyType(reflect.String)
	bytesType   = keyType(reflect.Slice)
)

var numberLikeKinds = [...]bool{
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Uintptr: true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  false,
	reflect.Slice:   false,
}

func (kt keyType) Compare(a, b interface{}) int {
	val1 := reflect.ValueOf(a)
	val2 := reflect.ValueOf(b)
	kind := kt.kind()
	result := compareTypes(val1, val2, kind)

	return result
}

func (kt keyType) CalcScore(key interface{}) float64 {
	k := reflect.ValueOf(key)
	kind := kt.kind()

	if kk := k.Kind(); kk != kind {
		// 特殊情况处理
		//
		if numberLikeKinds[kind] && (kk == reflect.Int || kk == reflect.Float64) {

		} else {
			name := kind.String()

			if kind == reflect.Slice {
				name = "[]byte"
			}

			panic(fmt.Errorf("skiplist: key type must be %v, but actual type is %v", name, k.Type()))
		}
	}

	score := calcScore(k)

	return score
}

// compareTypes 根据类型进行比较
func compareTypes(lhs, rhs reflect.Value, kind reflect.Kind) int {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		v1 := calcScore(lhs)
		v2 := calcScore(rhs)

		if v1 > v2 {
			return 1
		}

		if v1 < v2 {
			return -1
		}

		return 0

	case reflect.Int64:
		v1 := lhs.Int()
		v2 := rhs.Int()

		if v1 > v2 {
			return 1
		}

		if v1 < v2 {
			return -1
		}

		return 0

	case reflect.Uint64:
		v1 := lhs.Uint()
		v2 := rhs.Uint()

		if v1 > v2 {
			return 1
		}

		if v1 < v2 {
			return -1
		}

		return 0

	case reflect.String:
		v1 := lhs.String()
		v2 := rhs.String()

		if v1 == v2 {
			return 0
		}

		if v1 > v2 {
			return 1
		}

		return -1

	case reflect.Slice:
		if lhs.Type().ConvertibleTo(typeOfBytes) && rhs.Type().ConvertibleTo(typeOfBytes) {
			bytes1 := lhs.Convert(typeOfBytes).Interface().([]byte)
			bytes2 := rhs.Convert(typeOfBytes).Interface().([]byte)
			return bytes.Compare(bytes1, bytes2)
		}
	}
	return 0
}

// calcScore 计算key的分数
func calcScore(val reflect.Value) (score float64) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		score = float64(val.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		score = float64(val.Uint())

	case reflect.Float32, reflect.Float64:
		score = val.Float()

	case reflect.String:
		// 将字符串进行哈希计算
		var hash uint64
		str := val.String()
		l := len(str)

		// 只计算前8个字节
		if l > 8 {
			l = 8
		}

		// 将字符串转换为uint64
		for i := 0; i < l; i++ {
			shift := uint(64 - 8 - i*8)
			hash |= uint64(str[i]) << shift
		}

		score = float64(hash)

	case reflect.Slice:
		// 将[]byte转换为uint64
		if val.Type().ConvertibleTo(typeOfBytes) {
			var hash uint64
			data := val.Convert(typeOfBytes).Interface().([]byte)

			l := len(data)

			// only use first 8 bytes
			if l > 8 {
				l = 8
			}

			// Consider str as a Big-Endian uint64.
			for i := 0; i < l; i++ {
				shift := uint(64 - 8 - i*8)
				hash |= uint64(data[i]) << shift
			}

			score = float64(hash)
		}
	}

	return
}

// kind 根据keyType获取kind
func (kt keyType) kind() reflect.Kind {
	return reflect.Kind(kt)
}
