package api_matchers

import (
	"fmt"

	"github.com/go-web-kits/testx"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func ResponseSuccess() types.GomegaMatcher {
	return &ResponseSuccessfullyMatcher{}
}

type ResponseSuccessfullyMatcher struct{}

func (matcher *ResponseSuccessfullyMatcher) Match(actual interface{}) (success bool, err error) {
	r := actual.(testx.RR)
	result, ok := r.ResponseBody["result"].(map[string]interface{})
	if result == nil || !ok {
		return false, fmt.Errorf(
			"Expected the request `%s` to succeed, but it didn't, response body:\n%s",
			r.API, r.ResponseBody,
		)
	}
	code := result["code"]
	// TODO: configurable
	// FIXME: unmarshal num
	if codeNum, ok := code.(float64); ok && codeNum == 0.0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (matcher *ResponseSuccessfullyMatcher) FailureMessage(actual interface{}) (message string) {
	r := actual.(testx.RR)
	return format.Message(
		fmt.Sprintf("the response of the request `%s`:\n%v", r.API, r),
		"to have the success code", 0,
	)
}

func (matcher *ResponseSuccessfullyMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	r := actual.(testx.RR)
	return format.Message(
		fmt.Sprintf("the response of the request `%s`:\n%v", r.API, r),
		"not to have the success code", 0,
	)
}
