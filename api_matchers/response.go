package api_matchers

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-web-kits/lab/business_error"
	"github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/matchersx"
	"github.com/go-web-kits/utils/mapx"
	"github.com/go-web-kits/utils/reflectx"
	"github.com/go-web-kits/utils/structx"
	"github.com/k0kubun/pp"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func ResponseData(data interface{}) types.GomegaMatcher {
	return Response(map[string]interface{}{
		// "result": map[string]interface{}{"code": 0, "message": "success"},
		"data": data,
	})
}

func ResponseList(list []interface{}) types.GomegaMatcher {
	return ResponseData(map[string]interface{}{"list": list})
}

func Response(expected interface{}) types.GomegaMatcher {
	switch exp := expected.(type) {
	case int:
		return &ResponseMatcher{HttpCode: exp}
	case error:
		return &ResponseMatcher{Err: exp}
	case map[string]interface{}, []interface{}, []map[string]interface{}:
		return &ResponseMatcher{Body: exp}
	default:
		if reflect.TypeOf(expected).Kind() == reflect.Struct {
			return &ResponseMatcher{Body: structx.ToJsonizeMap(expected)}
		}
		panic("this type of expected arg is not supported")
	}
}

type ResponseMatcher struct {
	HttpCode int
	Body     interface{}
	Err      error
}

func (matcher *ResponseMatcher) Match(actual interface{}) (success bool, err error) {
	r := actual.(testx.RR)
	// pp.Println(r)
	if matcher.Err != nil {
		return judgeErr(r, matcher.Err)
	} else if matcher.HttpCode != 0 {
		return judgeHttpCode(r, matcher.HttpCode)
	} else {
		return judgeBody(r, matcher.Body)
	}
}

func (matcher *ResponseMatcher) FailureMessage(actual interface{}) (message string) {
	return rMsg(matcher, actual, "")
}

func (matcher *ResponseMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return rMsg(matcher, actual, "not ")
}

func rMsg(matcher *ResponseMatcher, actual interface{}, negated string) string {
	r := actual.(testx.RR)
	if os.Getenv("LOG") != "" {
		_, _ = pp.Println(r)
	}

	if matcher.Err != nil {
		return format.Message("response error:\n"+format.Object(r.ResponseBody, 1), negated+"to be the error", format.Object(matcher.Err, 1))
	} else if matcher.HttpCode != 0 {
		return format.Message("response HTTP status -> "+fmt.Sprint(r.ResponseCode)+" <-", negated+"to be", matcher.HttpCode)
	} else {
		var body interface{}
		body = r.ResponseBody
		if reflectx.IsZero(body) {
			body = r.ResponseBodySlice
		}
		return format.Message("response body:\n"+format.Object(body, 1), negated+"to be", format.Object(matcher.Body, 1))
	}
}

// ========

func judgeHttpCode(r testx.RR, expectedCode int) (bool, error) {
	return r.ResponseCode == expectedCode, nil
}

func judgeErr(r testx.RR, expectedErr error) (bool, error) {
	bodyError := fmt.Errorf(
		"Expected the request `%s` to response error\n\n  `%v`, but it didn't, actual response is:\n\n  `%v`",
		r.API, expectedErr, r.ResponseBody,
	)

	switch e := expectedErr.(type) {
	case business_error.Renderable:
		var code interface{}
		if len(e.RenderCodePath()) > 1 {
			root, ok := r.ResponseBody[e.RenderCodePath()[0]].(map[string]interface{})
			if root == nil || !ok {
				return false, bodyError
			}
			code = mapx.Dig(root, e.RenderCodePath()[1:]...)
		} else {
			code = r.ResponseBody[e.RenderCodePath()[0]]
		}

		if codeNum, ok := code.(float64); ok && codeNum == float64(e.GetCode()) {
			return true, nil
		}
		return false, nil
	default:
		var message string
		if result, ok := r.ResponseBody["result"].(map[string]interface{}); ok {
			message = result["message"].(string)
		} else if r.ResponseBody["msg"] != nil {
			message = r.ResponseBody["msg"].(string)
		} else if r.ResponseBody["message"] != nil {
			message = r.ResponseBody["message"].(string)
		} else {
			return false, bodyError
		}

		if !strings.Contains(message, e.Error()) {
			return false, nil
		}
		return true, nil
	}
}

func judgeBody(r testx.RR, expected interface{}) (bool, error) {
	switch exp := expected.(type) {
	case []interface{}, []map[string]interface{}:
		return matchersx.JudgeLike(r.ResponseBodySlice, exp), nil
	case map[string]interface{}:
		if r.ResponseBody == nil {
			return false, nil
		}
		keys := mapx.Keys(exp)
		onlyCheckData := len(keys) == 1 && keys[0] == "data"
		if onlyCheckData {
			return matchersx.JudgeLike(r.ResponseBody["data"], exp["data"]), nil
		}

		for k, expVal := range exp {
			if !matchersx.JudgeLike(r.ResponseBody[k], expVal) {
				return false, nil
			}
		}

		return true, nil
	}
	return false, nil
}
