package testx

import (
	"github.com/go-web-kits/utils/mapx"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var TestSubject interface{}

func IsExpected(params ...map[string]interface{}) AssertionX {
	if TestSubject != nil {
		return AssertionX{gomega.Expect(TestSubject), TestSubject}
	}

	if CurrentAPI != "" {
		if len(params) > 0 {
			ps := map[string]interface{}{}
			for _, p := range params {
				ps = mapx.Merge(ps, p)
			}

			return ExpectRequestedBy(ps)
		} else {
			return ExpectRequested()
		}
	}

	ginkgo.Fail("`TestSubject` is not set")
	return AssertionX{}
}
