package slog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AsString convers val into string.
func AsString(val interface{}) (string, bool) {
	if val == nil {
		return "", true
	}

	if str, ok := val.(string); ok {
		return str, true
	} else if strp, ok := val.(*string); ok {
		return *strp, true
	} else if tm, ok := val.(time.Time); ok {
		return tm.Format(time.RFC3339), true
	} else if tmp, ok := val.(*time.Time); ok {
		return tmp.Format(time.RFC3339), true
	} else if istr, ok := val.(interface{ String() string }); ok {
		return istr.String(), true
	} else {
		// format using simple Sprint
		return fmt.Sprint(val), false
	}
}

// AsStringQ convert value into string.
// Value will be double-qouted if value is string/time/stringer
func AsStringQ(val interface{}) string {
	str, isStr := AsString(val)
	if isStr {
		return strconv.Quote(str)
	}
	return str
}

// Simpleformatter return simple key-value formatter, separated with `sep`
func SimpleFormatter(msg string, keyVals []interface{}, sep string) string {
	sb := strings.Builder{}
	sb.WriteString(msg)
	n := len(keyVals)
	if n == 0 {
		return sb.String()
	}

	// number of pair
	nkv := (n + 1) / 2

	// number of args is odd
	if n%2 != 0 {
		n--
	}

	// format field and values
	j := 0
	for i := 0; i < n; i += 2 {
		field, _ := AsString(keyVals[i])
		sb.WriteString(field)
		sb.WriteString(sep)
		sb.WriteString(AsStringQ(keyVals[i+1]))
		j++
	}

	// number of args is odd
	if n != len(keyVals) {
		sb.WriteString(fmt.Sprintf("%s-%02d", UnknownFieldName, nkv))
		sb.WriteString(sep)
		sb.WriteString(AsStringQ(keyVals[n]))
	}

	return sb.String()
}

// FieldsToMap convert key-value array to map[string]interface{}
func FieldsToMap(keyVals []interface{}) map[string]interface{} {
	n := len(keyVals)
	if n == 0 {
		return nil
	}

	nkv := (n + 1) / 2
	kvMaps := make(map[string]interface{})

	// number of args is odd
	if n%2 != 0 {
		n--
	}

	// extract fields and values
	for i := 0; i < n; i += 2 {
		field := keyVals[i]
		var key string
		if field != nil {
			key, _ = AsString(field)
		} else {
			key = fmt.Sprintf("%s-%02d", UnknownFieldName, (i)/2+1)
		}
		kvMaps[key] = keyVals[i+1]
	}

	// number of args is odd
	if n != len(keyVals) {
		key := fmt.Sprintf("%s-%02d", UnknownFieldName, nkv)
		kvMaps[key] = keyVals[n]
	}

	return kvMaps
}

// SeparateFields into array of keys and array of values
func SeparateFields(keyVals []interface{}) ([]string, []interface{}) {
	n := len(keyVals)
	if n == 0 {
		return nil, nil
	}

	nkv := (n + 1) / 2
	fields := make([]string, nkv)
	values := make([]interface{}, nkv)

	// number of args is odd
	if n%2 != 0 {
		n--
	}

	// extract fields and values
	j := 0
	for i := 0; i < n; i += 2 {
		field := keyVals[i]
		if field != nil {
			str, _ := AsString(field)
			fields[j] = str
		} else {
			fields[j] = fmt.Sprintf("%s-%02d", UnknownFieldName, j+1)
		}

		values[j] = keyVals[i+1]
		j++
	}

	// number of args is odd
	if n != len(keyVals) {
		fields[nkv-1] = fmt.Sprintf("%s-%02d", UnknownFieldName, nkv)
		values[nkv-1] = keyVals[n]
	}

	return fields, values
}
