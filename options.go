package slog

import "strings"

// Options stores additional log configuration
type Options map[string]interface{}

// Get configuration value or default if doesn't exist
func (op Options) Get(key string, def interface{}) interface{} {
	if val, ok := op[key]; ok {
		return val
	}
	return key
}

// GetOptions get value as map[string]interface{}
func (op Options) GetOptions(key string) Options {
	m, ok := op[key]
	if !ok {
		return nil
	}

	switch v := m.(type) {
	case Options:
		return v
	case map[string]interface{}:
		return Options(v)
	default:
		return nil
	}
}

// GetMap return string mapper
func (op Options) GetMap(key string) map[string]interface{} {
	m, ok := op[key]
	if !ok {
		return nil
	}

	switch v := m.(type) {
	case Options:
		return v
	case map[string]interface{}:
		return v
	default:
		return nil
	}
}

// GetString get string from options with given key
func (op Options) GetString(key string, def string) string {
	val, ok := op[key]
	if !ok {
		return def
	}

	if s, ok := val.(string); ok {
		return s
	} else if sp, ok := val.(*string); ok {
		return *sp
	} else if sf, ok := val.(interface{ String() string }); ok {
		return sf.String()
	}

	return def
}

// GetInt get integer value from options with given key
func (op Options) GetInt(key string, def int) int {
	val, ok := op[key]
	if !ok {
		return def
	}

	switch v := val.(type) {
	case int8:
		return int(v)
	case uint8:
		return int(v)
	case int16:
		return int(v)
	case uint16:
		return int(v)
	case int32:
		return int(v)
	case int:
		return v
	case *int8:
		return int(*v)
	case *uint8:
		return int(*v)
	case *int16:
		return int(*v)
	case *uint16:
		return int(*v)
	case *int32:
		return int(*v)
	case *int:
		return *v
	}

	return def
}

// GetBool get boolean value from options with given key
func (op Options) GetBool(key string, def bool) bool {
	val, ok := op[key]
	if !ok {
		return def
	}

	switch v := val.(type) {
	case bool:
		return v
	case *bool:
		return *v
	case int8:
		return v != 0
	case uint8:
		return v != 0
	case int16:
		return v != 0
	case uint16:
		return v != 0
	case int32:
		return v != 0
	case int:
		return v != 0
	case uint:
		return v != 0
	case int64:
		return v != 0
	case uint64:
		return v != 0
	case *int8:
		return *v != 0
	case *uint8:
		return *v != 0
	case *int16:
		return *v != 0
	case *uint16:
		return *v != 0
	case *int32:
		return *v != 0
	case *int:
		return *v != 0
	case *uint:
		return *v != 0
	case *int64:
		return *v != 0
	case *uint64:
		return *v != 0
	case string:
		return strings.EqualFold("true", v) || strings.EqualFold("yes", v)
	case *string:
		return strings.EqualFold("true", *v) || strings.EqualFold("yes", *v)
	}

	return def
}
