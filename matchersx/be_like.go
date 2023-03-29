package matchersx

import (
	"github.com/onsi/gomega/format"
)

type BeLikeMatcher struct {
	Expected interface{}
}

func (matcher *BeLikeMatcher) Match(actual interface{}) (success bool, err error) {
	return JudgeLike(actual, matcher.Expected), nil
}

func (matcher *BeLikeMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to be like", matcher.Expected)
}

func (matcher *BeLikeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to be like", matcher.Expected)
}
