package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

/*
 * 说明：
 * 作者：吕元龙
 * 时间 2024/8/30 21:08
 */

// IntToBytes 将整数转换为字节切片
func IntToBytes(n int) []byte {
	// 根据整数的大小，选择合适的字节切片长度
	var buf bytes.Buffer

	switch {
	case n >= -128 && n <= 127:
		binary.Write(&buf, binary.BigEndian, int8(n))
	case n >= -32768 && n <= 32767:
		binary.Write(&buf, binary.BigEndian, int16(n))
	case n >= -2147483648 && n <= 2147483647:
		binary.Write(&buf, binary.BigEndian, int32(n))
	default:
		binary.Write(&buf, binary.BigEndian, int64(n))
	}

	return buf.Bytes()
}

// FormatBytesAsBinary 将字节切片转换为二进制字符串切片
func FormatBytesAsBinary(data []byte) []string {
	formatted := make([]string, len(data))
	for i, b := range data {
		formatted[i] = fmt.Sprintf("%08b", b)
	}
	return formatted
}

// ShiftBytes 将字节数组中的二进制数据左移或者右移
// flag: 0 左移 1 右移
// bit: 移位位数
func ShiftBytes(data []byte, flag, bit int) ([]byte, error) {
	if flag != 0 && flag != 1 {
		return data, errors.New("参数错误")
	}
	if flag == 1 {
		return shiftBytesRight(data, bit), nil
	} else {
		return shiftBytesLeft(data, bit), nil
	}
}

// ShiftBytesRight 对[]byte进行右移n位的操作
func shiftBytesRight(data []byte, n int) []byte {
	byteShift := n / 8 // 计算需要移动多少个字节
	bitShift := n % 8  // 剩余的位移量
	result := make([]byte, len(data))

	for i := len(data) - 1; i >= 0; i-- {
		if i-byteShift-1 >= 0 && bitShift > 0 {
			result[i] = data[i-byteShift] >> bitShift
			result[i] |= data[i-byteShift-1] << (8 - bitShift)
		} else if i-byteShift >= 0 {
			result[i] = data[i-byteShift] >> bitShift
		}
	}
	return result
}

// ShiftBytesLeft 对[]byte进行左移n位的操作
func shiftBytesLeft(data []byte, n int) []byte {
	byteShift := n / 8 // 计算需要移动多少个字节
	bitShift := n % 8  // 剩余的位移量
	result := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		if i+byteShift+1 < len(data) && bitShift > 0 {
			result[i] = data[i+byteShift] << bitShift
			result[i] |= data[i+byteShift+1] >> (8 - bitShift)
		} else if i+byteShift < len(data) {
			result[i] = data[i+byteShift] << bitShift
		}
	}

	return result
}
