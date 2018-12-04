package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type Assert struct {
	*testing.T
}

func (a *Assert) Nil(actual interface{}, logs ...interface{}) {
	if !IsNil(actual) {
		logCaller()
		fmt.Printf("expect nil, got %#v. %s\n", actual, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) NotNil(actual interface{}, logs ...interface{}) {
	if IsNil(actual) {
		logCaller()
		fmt.Println("expect value not to be nil.")
		a.FailNow()
	}
}

func (a *Assert) True(actual interface{}, logs ...interface{}) {
	if !IsBoolMatch(true, actual) {
		logCaller()
		fmt.Printf("expect value is not match true %#v. %s\n", actual, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) False(actual interface{}, logs ...interface{}) {
	if !IsBoolMatch(false, actual) {
		logCaller()
		fmt.Printf("expect value is not match false %#v. %s\n", actual, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) Bool(expect bool, actual interface{}, logs ...interface{}) {
	if !IsBoolMatch(expect, actual) {
		logCaller()
		fmt.Printf("expect boolean value is not equal %#v. %s\n", actual, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) Equal(expect, actual interface{}, logs ...interface{}) {
	if !IsEqual(expect, actual) {
		logCaller()
		fmt.Printf("expect value is not equal %#v. %s\n", actual, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) NotEqual(expect, actual interface{}, logs ...interface{}) {
	if IsEqual(expect, actual) {
		logCaller()
		fmt.Println("expect value is equal to actual.")
		a.FailNow()
	}
}

func (a *Assert) Contains(list, element interface{}, logs ...interface{}) {
	if ok, found := IncludeElement(list, element); !ok || !found {
		logCaller()
		fmt.Printf("expect list is not contains %#v. %s\n", element, fmt.Sprint(logs...))
		a.FailNow()
	}
}

func (a *Assert) NotContains(list, element interface{}, logs ...interface{}) {
	if ok, found := IncludeElement(list, element); !ok || found {
		logCaller()
		fmt.Printf("expect list is contains %#v. %s\n", element, fmt.Sprint(logs...))
		a.FailNow()
	}
}

// New return a pointer of new Assert struct value
func New(t *testing.T) *Assert {
	if nil == t {
		return nil
	}
	return &Assert{t}
}

// IsNil checks if a specified object is nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	case reflect.Bool, reflect.UnsafePointer,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.Array, reflect.String, reflect.Struct:
		// no nil value for these types
		return false
	}
	return false
}

// IsBoolMatch checks if the actual value is semantics match expected bool
func IsBoolMatch(expect bool, actual interface{}) bool {
	value := reflect.ValueOf(actual)
	switch value.Kind() {
	case reflect.Bool:
		if expect == value.Bool() {
			return true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := value.Int()
		b := i == 0 || i == 1
		var e int64
		if expect {
			e = 1
		}
		if b && e == i {
			return true
		}
	case reflect.String:
		b, err := strconv.ParseBool(value.String())
		if nil == err && expect == b {
			return true
		}
	}
	return IsEqual(expect, actual)
}

// IsEqual check if two values are equal
func IsEqual(expect, actual interface{}) bool {
	if IsFunction(expect) || IsFunction(actual) {
		return false
	}
	if expect == nil || actual == nil {
		return expect == actual
	}

	var e, a []byte
	var ok bool
	if e, ok = expect.([]byte); !ok {
		return reflect.DeepEqual(expect, actual)
	}
	if a, ok = actual.([]byte); !ok {
		return false
	}
	if e == nil || a == nil {
		return e == nil && a == nil
	}
	return bytes.Equal(e, a)
}

// IsFunction check if v is a function
func IsFunction(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Func
}

// IncludeElement check if the list contains the given element
func IncludeElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if IsEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if IsEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

// logCaller print the file path and current running line number of function caller
func logCaller(n ...int) {
	skip := 2
	if len(n) == 1 {
		skip = n[0]
	}
	_, file, line, _ := runtime.Caller(skip)
	fmt.Printf("%s:%d: ", file, line)
}
