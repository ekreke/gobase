package maputil

import (
	"slices"
	"testing"
)

type testStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestStructToMapE_Success(t *testing.T) {
	s := testStruct{A: 1, B: "x"}
	m, err := StructToMapE(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 keys in map, got %d", len(m))
	}
	if v, ok := m["a"]; !ok {
		t.Fatalf("expected key %q in map", "a")
	} else {
		if num, ok := v.(float64); !ok || num != 1 {
			t.Fatalf("expected key %q value 1, got %#v", "a", v)
		}
	}
	if v, ok := m["b"]; !ok {
		t.Fatalf("expected key %q in map", "b")
	} else if v != "x" {
		t.Fatalf("expected key %q value %q, got %#v", "b", "x", v)
	}
}

func TestStructToMapE_Error(t *testing.T) {
	ch := make(chan int)
	m, err := StructToMapE(ch)
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil, map=%v", m)
	}
}

func TestStructToMap_NoErrorVersion(t *testing.T) {
	ch := make(chan int)
	m := StructToMap(ch)
	if m == nil {
		t.Fatalf("expected non-nil map on error, got nil")
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map on error, got %v", m)
	}

	s := testStruct{A: 2, B: "y"}
	m = StructToMap(s)
	if v, ok := m["b"]; !ok || v != "y" {
		t.Fatalf("expected key %q value %q, got %v", "b", "y", m["b"])
	}
}

func TestMapToStructE_Success(t *testing.T) {
	m := map[string]interface{}{"a": 3, "b": "z"}
	s, err := MapToStructE[testStruct](m)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.A != 3 || s.B != "z" {
		t.Fatalf("expected struct {A:3 B:\"z\"}, got %+v", s)
	}
}

func TestMapToStructE_Error(t *testing.T) {
	m := map[string]interface{}{"ch": make(chan int)}
	_, err := MapToStructE[testStruct](m)
	if err == nil {
		t.Fatalf("expected error for unsupported value in map, got nil")
	}
}

func TestMapToStruct_NoErrorVersion(t *testing.T) {
	m := map[string]interface{}{"a": 5, "b": "w"}
	s := MapToStruct[testStruct](m)
	if s.A != 5 || s.B != "w" {
		t.Fatalf("expected struct {A:5 B:\"w\"}, got %+v", s)
	}

	bad := map[string]interface{}{"ch": make(chan int)}
	s = MapToStruct[testStruct](bad)
	if s.A != 0 || s.B != "" {
		t.Fatalf("expected zero value on error, got %+v", s)
	}
}

func TestToStringE_Success(t *testing.T) {
	m := map[string]interface{}{"a": 1, "b": "x"}
	s, err := ToStringE(m)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s == "" {
		t.Fatalf("expected non-empty json string")
	}
	if s[0] != '{' || s[len(s)-1] != '}' {
		t.Fatalf("expected json object, got %q", s)
	}
}

func TestToStringE_Error(t *testing.T) {
	m := map[string]interface{}{"ch": make(chan int)}
	s, err := ToStringE(m)
	if err == nil {
		t.Fatalf("expected error, got nil with result %q", s)
	}
}

func TestToString_UsesDefaultOnError(t *testing.T) {
	m := map[string]interface{}{"ch": make(chan int)}
	s := ToString(m)
	if s != "{}" {
		t.Fatalf("expected default \"{}\" on error, got %q", s)
	}
}

func TestGetMapValue(t *testing.T) {
	m := map[string]int{"a": 1}
	if v := GetMapValue(m, "a", 0); v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
	if v := GetMapValue(m, "b", 5); v != 5 {
		t.Fatalf("expected default 5, got %d", v)
	}
}

func TestMergeMap(t *testing.T) {
	base := map[string]int{"a": 1, "b": 2}
	m1 := map[string]int{"b": 3}
	m2 := map[string]int{"c": 4}
	got := MergeMap(base, m1, m2)
	if len(got) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(got))
	}
	if got["a"] != 1 || got["b"] != 3 || got["c"] != 4 {
		t.Fatalf("unexpected merged map: %#v", got)
	}
}

func TestMapDecode(t *testing.T) {
	m := map[string]string{
		"a": "hello%20world",
		"b": "%zz",
	}
	got := MapDecode(m)
	if got["a"] != "hello world" {
		t.Fatalf("expected decoded value, got %q", got["a"])
	}
	if got["b"] != "%zz" {
		t.Fatalf("expected original value on decode error, got %q", got["b"])
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	vals := MapValues(m)
	if len(vals) != 3 {
		t.Fatalf("expected 3 values, got %d", len(vals))
	}
	slices.Sort(vals)
	if !slices.Equal(vals, []int{1, 2, 3}) {
		t.Fatalf("unexpected values slice: %#v", vals)
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	slices.Sort(keys)
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected keys slice: %#v", keys)
	}
}
