package matchersx

import (
	"bytes"
	"reflect"

	"github.com/go-web-kits/utils/slicex"
	"github.com/go-web-kits/utils/structx"
	"github.com/onsi/gomega/format"
)

type HaveAttributesMatcher struct {
	Expected interface{}
	Ignore   []string
}

type IncludeMatcher = HaveAttributesMatcher

func (matcher *HaveAttributesMatcher) Match(actual interface{}) (success bool, err error) {
	expectedMap := structx.ToJsonizeMap(matcher.Expected)
	actualMap := structx.ToJsonizeMap(actual)
	givenStruct := reflect.TypeOf(matcher.Expected).Kind() == reflect.Struct
	return include(actualMap, expectedMap, matcher.Ignore, givenStruct), nil
}

func (matcher *HaveAttributesMatcher) FailureMessage(actual interface{}) (message string) {
	expectedMap := structx.ToJsonizeMap(matcher.Expected)
	actualMap := structx.ToJsonizeMap(actual)
	return format.Message(actualMap, "to have attributes", expectedMap)
}

func (matcher *HaveAttributesMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	expectedMap := structx.ToJsonizeMap(matcher.Expected)
	actualMap := structx.ToJsonizeMap(actual)
	return format.Message(actualMap, "not to have attributes", expectedMap)
}

// ==========

// actual judgeMapInclude expected
func include(actual, expected map[string]interface{}, ignore []string, givenStruct bool) bool {
	for k, expVal := range expected {
		if slicex.IncludeStr(append(ignore, []string{"created_at", "updated_at"}...), k) {
			continue // TODO
		}

		if givenStruct && expVal != nil && expVal == reflect.Zero(reflect.TypeOf(expVal)).Interface() {
			continue
		}

		if exp, ok := expVal.(map[string]interface{}); ok {
			if act, ok := actual[k].(map[string]interface{}); ok {
				if include(act, exp, ignore, givenStruct) {
					continue
				} else {
					return false
				}
			}
		}

		if !equal(actual[k], expVal) {
			return false
		}
	}

	return true
}

func equal(actual interface{}, expected interface{}) bool {
	if actual == nil && expected == nil {
		return true //false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	} else if actual == nil || expected == nil {
		return false
	}
	// Shortcut for byte slices.
	// Comparing long byte slices with reflect.DeepEqual is very slow,
	// so use bytes.Equal if actual and expected are both byte slices.
	if actualByteSlice, ok := actual.([]byte); ok {
		if expectedByteSlice, ok := expected.([]byte); ok {
			return bytes.Equal(actualByteSlice, expectedByteSlice)
		}
	}

	convertedActual := actual
	if reflect.TypeOf(actual).ConvertibleTo(reflect.TypeOf(expected)) {
		convertedActual = reflect.ValueOf(actual).Convert(reflect.TypeOf(expected)).Interface()
	}
	return reflect.DeepEqual(convertedActual, expected)
}
