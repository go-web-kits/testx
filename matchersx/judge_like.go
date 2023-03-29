package matchersx

import (
	"bytes"
	"reflect"
)

func JudgeLike(actual interface{}, expected interface{}) bool {
	if actual == nil && expected == nil {
		return true // TODO, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	} else if actual == nil || expected == nil {
		return false
	}

	// Shortcut for byte slices.
	// Comparing long byte slices with reflect.DeepEqual is very slow,
	// so use bytes.Equal if actual and expected are both byte slices.
	if actualByteSlice, ok := actual.([]byte); ok {
		if expectedByteSlice, ok := expected.([]byte); ok {
			return bytes.Equal(actualByteSlice, expectedByteSlice) //, nil
		}
	}

	actVal, expVal := reflect.Indirect(reflect.ValueOf(actual)), reflect.Indirect(reflect.ValueOf(expected))
	actTp, expTp := actVal.Type(), expVal.Type()

	if expTp.Kind() == reflect.Map && actTp.Kind() == reflect.Map {
		// actKeys := actVal.MapKeys()
		for _, key := range expVal.MapKeys() {
			if !JudgeLike(actVal.MapIndex(key).Interface(), expVal.MapIndex(key).Interface()) {
				return false
			}
		}
		return true
	}

	if expTp.Kind() == reflect.Slice && actTp.Kind() == reflect.Slice {
		if actVal.Len() != expVal.Len() {
			return false
		}

		for i := 0; i < expVal.Len(); i++ {
			foundActIndex := foundIndexInSlice(expVal.Index(i).Interface(), actVal)
			if foundActIndex == -1 {
				return false
			}

			actVal = reflect.AppendSlice(
				actVal.Slice(0, foundActIndex),
				actVal.Slice(foundActIndex+1, actVal.Len()))
		}

		return true
	}

	convertedActual := actual
	if actTp.ConvertibleTo(expTp) {
		convertedActual = actVal.Convert(expTp).Interface()
	}
	return reflect.DeepEqual(convertedActual, expected) //, nil
}

func foundIndexInSlice(expItem interface{}, actual reflect.Value) int {
	for i := 0; i < actual.Len(); i++ {
		if JudgeLike(actual.Index(i).Interface(), expItem) {
			return i
		}
	}
	return -1
}
