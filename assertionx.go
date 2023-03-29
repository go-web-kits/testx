package testx

import "github.com/onsi/gomega"

type AssertionX struct {
	gomega.Assertion
	Subject interface{}
}

func Expectx(subject interface{}) AssertionX {
	return AssertionX{gomega.Expect(subject), subject}
}

func (a AssertionX) ResponseBody() AssertionX {
	body := a.Subject.(RR).ResponseBody
	return AssertionX{gomega.Expect(body), body}
}

func (a AssertionX) ResponseData() AssertionX {
	data := a.Subject.(RR).ResponseBody["data"]
	return AssertionX{gomega.Expect(data), data}
}

func (a AssertionX) ResponseCode() AssertionX {
	code := a.Subject.(RR).ResponseCode
	return AssertionX{gomega.Expect(code), code}
}
