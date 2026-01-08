package osutil

import (
	"errors"
	"fmt"
	"reflect"
)

// Call invokes the method with the given parameters.
// It returns the values returned by the method and an error if any.
func Call(method interface{}, params ...interface{}) ([]reflect.Value, error) {
	t := reflect.TypeOf(method)
	if t.Kind() != reflect.Func {
		return nil, errors.New("the input is not a function")
	}

	f := reflect.ValueOf(method)

	if len(params) != t.NumIn() {
		return nil, errors.New("the number of input params not match")
	}

	// convert params to reflect.Value
	var mp = make([]reflect.Value, len(params))
	for i, v := range params {
		inType := t.In(i)
		val := reflect.ValueOf(v)

		// check the type of param
		if v == nil {
			if inType.Kind() != reflect.Pointer && inType.Kind() != reflect.Interface && inType.Kind() != reflect.Slice && inType.Kind() != reflect.Map && inType.Kind() != reflect.Chan && inType.Kind() != reflect.Func {
				return nil, fmt.Errorf("param[%d] cannot be nil for type %s", i, inType)
			}
			mp[i] = reflect.Zero(inType)
		} else {
			if !val.Type().AssignableTo(inType) {
				return nil, fmt.Errorf("param[%d] type mismatch: expected %s, got %s", i, inType, val.Type())
			}
			mp[i] = val
		}
		mp[i] = val
	}
	return f.Call(mp), nil
}
