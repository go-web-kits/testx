package matchersx

import (
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/utils"
	"github.com/onsi/gomega/format"
)

type AffectMatcher struct{}

func (matcher *AffectMatcher) Match(actual interface{}) (success bool, err error) {
	result := actual.(dbx.Result)
	if result.Err == nil {
		return true, nil
	}

	return false, nil
}

func (matcher *AffectMatcher) FailureMessage(actual interface{}) (message string) {
	result := actual.(dbx.Result)
	return format.Message("database operation on "+utils.TypeNameOf(result.Data),
		"not to have error", format.Object(result.Err, 1))
}

func (matcher *AffectMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	result := actual.(dbx.Result)
	return format.Message("database operation on "+utils.TypeNameOf(result.Data), "to have error")
}
