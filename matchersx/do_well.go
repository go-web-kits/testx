package matchersx

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type DoWellMatcher struct{}

func (matcher *DoWellMatcher) Match(actual interface{}) (success bool, err error) {
	// is purely nil?
	if actual == nil {
		return true, nil
	}

	actualVal := reflect.ValueOf(actual)
	if actualVal.Kind() == reflect.Slice {
		actual = actualVal.Index(actualVal.Len() - 1).Interface()
	}

	// must be an 'error' type
	if !isError(actual) {
		return false, fmt.Errorf("Expected an error-type.  Got:\n%s", format.Object(actual, 1))
	}

	// must be nil (or a pointer to a nil)
	return isNil(actual), nil
}

func (matcher *DoWellMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected success, but got an error:\n%s\n%s", format.Object(actual, 1), format.IndentString(actual.(error).Error(), 1))
}

func (matcher *DoWellMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "Expected failure, but got no error."
}

// ========================

func isError(a interface{}) bool {
	_, ok := a.(error)
	return ok
}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(a).IsNil()
	}

	return false
}
