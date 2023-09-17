package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// MatchPredicate helper
func MatchPredicate(
	value []byte,
	predicate map[string]string,
) bool {
	if len(predicate) == 0 {
		return true
	}
	m := make(map[string]interface{})
	err := json.Unmarshal(value, &m)
	if err != nil {
		return false
	}
	for k, v := range predicate {
		op := "=="
		if strings.HasSuffix(k, ":==") {
			k = k[0 : len(k)-3]
			op = "=="
		} else if strings.HasSuffix(k, ":!=") {
			k = k[0 : len(k)-3]
			op = "!="
		} else if strings.HasSuffix(k, ":>=") {
			k = k[0 : len(k)-3]
			op = ">="
		} else if strings.HasSuffix(k, ":<=") {
			k = k[0 : len(k)-3]
			op = "<="
		} else if strings.HasSuffix(k, ":<") {
			k = k[0 : len(k)-2]
			op = "<"
		} else if strings.HasSuffix(k, ":>") {
			k = k[0 : len(k)-2]
			op = ">"
		}
		if m[k] == nil {
			return false
		}
		valType := reflect.TypeOf(m[k]).String()
		if strings.HasPrefix(valType, "int") ||
			strings.HasPrefix(valType, "uint") {
			val1 := ToInt64(m[k])
			val2 := ToInt64(v)
			if (op == "==" && val1 != val2) ||
				(op == "!=" && val1 == val2) ||
				(op == ">=" && val1 < val2) ||
				(op == "<=" && val1 > val2) ||
				(op == ">" && val1 <= val2) ||
				(op == "<" && val1 >= val2) {
				return false
			}
		} else if strings.HasPrefix(valType, "float") {
			val1 := ToFloat64(m[k])
			val2 := ToFloat64(v)
			if (op == "==" && val1 != val2) ||
				(op == "!=" && val1 == val2) ||
				(op == ">=" && val1 < val2) ||
				(op == "<=" && val1 > val2) ||
				(op == ">" && val1 <= val2) ||
				(op == "<" && val1 >= val2) {
				return false
			}
		} else {
			str := fmt.Sprintf("%v", m[k])
			if (op == "==" && str != v) ||
				(op == "!=" && str == v) ||
				(op == ">=" && str < v) ||
				(op == "<=" && str > v) ||
				(op == ">" && str <= v) ||
				(op == "<" && str >= v) {
				return false
			}
		}
	}
	return true
}
