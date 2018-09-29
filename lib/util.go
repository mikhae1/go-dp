package lib

import (
	"reflect"
)

// ArrayContains finds index of `val` in `array`
func ArrayContains(array interface{}, val interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				return
			}
		}
	}

	return
}

// ArrayContainsArray checks srcArray in dstArray
func ArrayContainsArray(srcArray interface{}, dstArray interface{}) (contains bool, srcIndex int) {
	srcIndex = -1

	if reflect.TypeOf(srcArray).Kind() != reflect.TypeOf(dstArray).Kind() &&
		reflect.TypeOf(srcArray).Kind() != reflect.Slice {
		return
	}

	srcVal := reflect.ValueOf(srcArray)
	dstVal := reflect.ValueOf(dstArray)

	for s := 0; s < srcVal.Len(); s++ {
		si := srcVal.Index(s).Interface()

		found := false
		for d := 0; d < dstVal.Len(); d++ {
			di := dstVal.Index(d).Interface()

			if reflect.DeepEqual(si, di) {
				found = true
				break
			}
		}

		if !found {
			srcIndex = s
			return
		}
	}

	contains = true
	return
}
