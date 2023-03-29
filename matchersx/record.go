package matchersx

import (
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/utils/slicex"
	"github.com/onsi/gomega/format"
)

type BeTheSameRecordMatcher struct {
	Expected interface{}
}

func (matcher *BeTheSameRecordMatcher) Match(actual interface{}) (success bool, err error) {
	if dbx.IdOf(actual) == dbx.IdOf(matcher.Expected) {
		return true, nil
	}

	return false, nil
}

func (matcher *BeTheSameRecordMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be the same record (have the same id) to", matcher.Expected)
}

func (matcher *BeTheSameRecordMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be the same record (have the same id) to", matcher.Expected)
}

type BeTheSameRecordsMatcher struct {
	Expected interface{}
}

func (matcher *BeTheSameRecordsMatcher) Match(actual interface{}) (success bool, err error) {
	actualIDs := (dbx.Result{Data: actual}).GetIds()
	expIDs := (dbx.Result{Data: matcher.Expected}).GetIds()
	if len(actualIDs) != len(expIDs) {
		return false, nil
	}

	for _, aid := range actualIDs {
		if !slicex.IncludeUint(expIDs, aid) {
			return false, nil
		}
	}

	return true, nil
}

func (matcher *BeTheSameRecordsMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be the same records (have the same ids) to", matcher.Expected)
}

func (matcher *BeTheSameRecordsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be the same records (have the same ids) to", matcher.Expected)
}
