package maputil

import (
	"maps"
	"net/url"

	"github.com/bytedance/sonic"
)

// GetMapValue get map value with default value
func GetMapValue[T comparable, A any](data map[T]A, property T, def A) A {
	if _, ok := data[property]; ok {
		return data[property]
	}
	return def
}

// MergeMap merge multiple map to one
func MergeMap[T comparable, A any](data map[T]A, d ...map[T]A) map[T]A {
	if len(d) > 0 {
		for _, val := range d {
			maps.Copy(data, val)
		}
	}
	return data
}

// MapDecode decode map value with url.QueryUnescape
func MapDecode[T comparable](data map[T]string) map[T]string {
	for kl, vl := range data {
		tmp, err := url.QueryUnescape(vl)
		if err != nil {
			tmp = vl
		}
		data[kl] = tmp
	}
	return data
}

// MapValues return map values as slice
func MapValues[T comparable, A any](data map[T]A) []A {
	var res = make([]A, 0)
	for _, v := range data {
		res = append(res, v)
	}
	return res
}

// MapKeys return map keys as slice
func MapKeys[T comparable, A any](data map[T]A) []T {
	var res = make([]T, 0)
	for k := range data {
		res = append(res, k)
	}
	return res
}

// ToString convert map to string
func ToString(m map[string]interface{}) string {
	s, err := ToStringE(m)
	if err != nil {
		return "{}"
	}
	return s
}

// ToStringE convert map to string with error handling
func ToStringE(m map[string]interface{}) (string, error) {
	bytes, err := sonic.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// StructToMap convert struct to map
func StructToMap(obj interface{}) map[string]interface{} {
	m, err := StructToMapE(obj)
	if err != nil {
		return map[string]interface{}{}
	}
	return m
}

// StructToMapE convert struct to map with error handling
func StructToMapE(obj interface{}) (map[string]interface{}, error) {
	jsonData, err := sonic.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = sonic.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MapToStruct convert map to struct
func MapToStruct[T any](m map[string]interface{}) T {
	out, _ := MapToStructE[T](m)
	return out
}

// MapToStructE convert map to struct with error handling
func MapToStructE[T any](m map[string]interface{}) (out T, err error) {
	jsonData, err := sonic.Marshal(m)
	if err != nil {
		return out, err
	}
	if err = sonic.Unmarshal(jsonData, &out); err != nil {
		return out, err
	}
	return out, nil
}
