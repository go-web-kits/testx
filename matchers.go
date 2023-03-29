package testx

import (
	"github.com/go-web-kits/testx/matchersx"
	"github.com/onsi/gomega/types"
)

func HaveAffected() types.GomegaMatcher {
	return &matchersx.AffectMatcher{}
}

func HaveFound(total ...int) types.GomegaMatcher {
	if len(total) == 0 {
		return &matchersx.HaveFoundMatcher{}
	} else {
		return &matchersx.HaveFoundMatcher{Total: total[0]}
	}
}

func BeTheSameRecordTo(expected interface{}) types.GomegaMatcher {
	return &matchersx.BeTheSameRecordMatcher{
		Expected: expected,
	}
}

func BeTheSameRecordsTo(expected ...interface{}) types.GomegaMatcher {
	return &matchersx.BeTheSameRecordsMatcher{
		Expected: expected,
	}
}
func Include(expected interface{}, ignore ...string) types.GomegaMatcher {
	return &matchersx.HaveAttributesMatcher{
		Expected: expected,
		Ignore:   ignore,
	}
}

var HaveAttributes = Include

func BeLike(expected interface{}) types.GomegaMatcher {
	return &matchersx.BeLikeMatcher{
		Expected: expected,
	}
}
