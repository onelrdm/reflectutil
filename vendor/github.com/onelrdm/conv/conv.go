// Package conv implements conversion functions between basic data types.
package conv

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrUnknownType = errors.New("unknown type")
)

// ShouldInt64 converts the value into int64.
func ShouldInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case fmt.Stringer:
		return strconv.ParseInt(v.String(), 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	}
	return 0, ErrUnknownType
}

// MustInt64 converts the value into int64.
// It panics if an error occurs in conversion.
func MustInt64(v interface{}) int64 {
	ret, err := ShouldInt64(v)
	if err != nil {
		panic(err)
	}
	return ret
}

// ShouldFloat64 converts the value into float64.
func ShouldFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case fmt.Stringer:
		return strconv.ParseFloat(v.String(), 64)
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, ErrUnknownType
}

// MustFloat64 converts the value into float64.
// It panics if an error occurs in conversion.
func MustFloat64(v interface{}) float64 {
	ret, err := ShouldFloat64(v)
	if err != nil {
		panic(err)
	}
	return ret
}

// ShouldString converts the value into string.
func ShouldString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case []byte:
		return string(v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case uint:
		return fmt.Sprintf("%d", v), nil
	case uint8:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int8:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return fmt.Sprintf("%f", v), nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	default:
		return "", ErrUnknownType
	}
}

// MustString converts the value into string.
// It panics if an error occurs in conversion.
func MustString(v interface{}) string {
	ret, err := ShouldString(v)
	if err != nil {
		panic(err)
	}
	return ret
}
