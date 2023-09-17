package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ArrayToMap helper
func ArrayToMap(args ...string) map[string]string {
	res := make(map[string]string)
	for i := 0; i < len(args)-1; i++ {
		res[args[i]] = args[i+1]
	}
	return res
}

// ToInt helper
func ToInt(input any) (res int) {
	if input == nil {
		return 0
	}
	switch input.(type) {
	case int:
		res = input.(int)
	case uint:
		res = int(input.(uint))
	case int32:
		res = int(input.(int32))
	case int64:
		res = int(input.(int64))
	default:
		res, _ = strconv.Atoi(fmt.Sprintf("%v", input))
	}
	return
}

// ToBoolean converter
func ToBoolean(input any) bool {
	if input == nil {
		return false
	}
	switch input.(type) {
	case int:
		return input.(int) != 0
	case uint:
		return input.(uint) != 0
	case int32:
		return input.(int32) != 0
	case int64:
		return input.(int64) != 0
	case float32:
		return input.(float32) != 0
	case float64:
		return input.(float64) != 0
	case bool:
		return input.(bool)
	case string:
		s := strings.ToLower(input.(string))
		return s == "true" || s == "yes" || s == "y"
	}
	return false
}

// ToFloat64 converter
func ToFloat64(input any) (res float64) {
	if input == nil {
		return 0
	}
	switch input.(type) {
	case int:
		res = float64(input.(int))
	case uint:
		res = float64(input.(uint))
	case int32:
		res = float64(input.(int32))
	case int64:
		res = float64(input.(int64))
	case float32:
		res = float64(input.(float32))
	case float64:
		res = input.(float64)
	case *float64:
		f := input.(*float64)
		if f != nil {
			res = *f
		}
	case *float32:
		f := input.(*float32)
		if f != nil {
			res = float64(*f)
		}
	case *uint64:
		f := input.(*uint64)
		if f != nil {
			res = float64(*f)
		}
	case *int64:
		f := input.(*int64)
		if f != nil {
			res = float64(*f)
		}
	default:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%v", input), 64)
	}
	return
}

// ToInt64 helper
func ToInt64(input any) (res int64) {
	if input == nil {
		return 0
	}
	switch input.(type) {
	case int:
		res = int64(input.(int))
	case uint:
		res = int64(input.(uint))
	case int32:
		res = int64(input.(int32))
	case int64:
		res = input.(int64)
	default:
		res, _ = strconv.ParseInt(fmt.Sprintf("%v", input), 10, 64)
	}
	return
}

// Includes helper
func Includes(arr []string, s string) bool {
	for _, next := range arr {
		if next == s {
			return true
		}
	}
	return false
}

// AddSlice adding element to slice
func AddSlice(arr []string, strings ...string) []string {
	for _, s := range strings {
		if !Includes(arr, s) {
			arr = append(arr, s)
		}

	}
	return arr
}

// RemoveSlice removing element from slice
func RemoveSlice(arr []string, strings ...string) []string {
	for _, s := range strings {
		for i, next := range arr {
			if next == s {
				arr = append(arr[:i], arr[i+1:]...)
			}
		}
	}
	return arr
}

// MapIntToArray helper
func MapIntToArray(m map[string]int32) (res []string) {
	for id := range m {
		res = append(res, id)
	}
	return
}
