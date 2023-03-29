package matchersx

import (
	"fmt"

	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/utils"
	"github.com/onsi/gomega/format"
)

type HaveFoundMatcher struct {
	Total int
}

func (matcher *HaveFoundMatcher) Match(actual interface{}) (success bool, err error) {
	result := actual.(dbx.Result)
	if result.NotFound() {
		return false, nil
	}

	if result.Err != nil {
		return false, fmt.Errorf(
			"Error occurred: \n%s",
			format.Object(result.Err, 1),
		)
	}

	if matcher.Total == 0 {
		return true, nil
	} else {
		return result.Total == matcher.Total, nil
	}
}

func (matcher *HaveFoundMatcher) FailureMessage(actual interface{}) (message string) {
	result := actual.(dbx.Result)
	return format.Message(utils.TypeNameOf(result.Data), "to be found")
}

func (matcher *HaveFoundMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	result := actual.(dbx.Result)
	return format.Message(utils.TypeNameOf(result.Data), "not to be found")
}
