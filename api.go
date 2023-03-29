package testx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"

	"github.com/go-web-kits/lab/routex"
	"github.com/go-web-kits/spear/bubo"
	"github.com/go-web-kits/utils"
	"github.com/go-web-kits/utils/mapx"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var CurrentAPI string
var CurrentParams map[string]interface{}

// =====
// Block
// =====

func API(subject interface{}, body func()) bool {
	name, ok := subject.(string)
	if !ok {
		name = utils.GetFuncName(subject)
	}
	r := ginkgo.Describe(name, func() {
		ginkgo.BeforeEach(func() {
			CurrentAPI = name
		})
		ginkgo.AfterEach(func() {
			CurrentAPI = ""
		})

		body()
	})
	return r
}

func SetParams(params map[string]interface{}) {
	ginkgo.BeforeEach(func() {
		CurrentParams = params
	})
	ginkgo.AfterEach(func() {
		CurrentParams = nil
	})
}

// =========
// Assertion
// =========

func ExpectRequested() AssertionX {
	subject := Request()
	return AssertionX{gomega.Expect(subject), subject}
}

func ExpectRequestedBy(params map[string]interface{}) AssertionX {
	subject := RequestBy(params)
	return AssertionX{gomega.Expect(subject), subject}
}

func ExpectRequestedWith(params map[string]interface{}) AssertionX {
	subject := RequestWith(params)
	return AssertionX{gomega.Expect(subject), subject}
}

// ================
// Request-Response
// ================

type RR struct {
	API                string
	Params             interface{}
	ResponseCode       int
	ResponseBody       map[string]interface{}
	ResponseBodySlice  []interface{}
	ResponseBodyString string
	ResponseHeader     http.Header
}

func Request() RR {
	return HTTPRequest("", "", mapx.Merge(CurrentParams, map[string]interface{}{}))
}

func RequestBy(params map[string]interface{}) RR {
	return HTTPRequest("", "", params)
}

func RequestWith(params map[string]interface{}) RR {
	return HTTPRequest("", "", mapx.Merge(CurrentParams, params))
}

func HTTPRequest(method, path string, param ...map[string]interface{}) RR {
	route, found := GetRoute(method, path)
	if !found {
		ginkgo.Fail("No Such Route: `" + path + "`")
		return RR{}
	}

	params := map[string]interface{}{}
	if len(param) > 0 {
		params = param[0]
	}

	fullPath := fillPathParams(route.Path, params)
	req := buildRequest(route.Method, fullPath, params)
	response := httptest.NewRecorder()
	CurrentApp.Engine.ServeHTTP(response, req)

	var body map[string]interface{}
	var bodySlice []interface{}
	if json.Unmarshal([]byte(response.Body.String()), &body) != nil {
		_ = json.Unmarshal([]byte(response.Body.String()), &bodySlice)
	}

	return RR{
		API:                path,
		ResponseCode:       response.Code,
		ResponseBody:       body,
		ResponseBodySlice:  bodySlice,
		ResponseBodyString: response.Body.String(),
		ResponseHeader:     response.Header(),
	}
}

func HTTPGet(path string, params ...map[string]interface{}) RR {
	return HTTPRequest("GET", path, params...)
}

func HTTPPost(path string, params ...map[string]interface{}) RR {
	return HTTPRequest("POST", path, params...)
}

func HTTPPut(path string, params ...map[string]interface{}) RR {
	return HTTPRequest("PUT", path, params...)
}

func HTTPDelete(path string, params ...map[string]interface{}) RR {
	return HTTPRequest("DElETE", path, params...)
}

// ===============
// private methods
// ===============

func toQuery(params map[string]interface{}) string {
	pairs := []string{}
	for k, v := range params {
		pairs = append(pairs, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(pairs, "&")
}

func fillPathParams(path string, params map[string]interface{}) string {
	pathParams := regexp.MustCompile(`/:[^/]*`).FindAllString(path, -1)
	if len(pathParams) > 0 {
		for _, p := range pathParams {
			value := params[p[2:]]
			if value == nil {
				panic("Invalid path parameter " + p[1:] + ", check your test code.")
			}
			path = strings.Replace(path, p, fmt.Sprintf("/%v", value), 1)
			delete(params, p[2:])
		}
	}
	return path
}

func buildRequest(method string, path string, params map[string]interface{}) *http.Request {
	var req *http.Request
	headers := params["headers"]
	delete(params, "headers")
	f := params["file"]
	delete(params, "file")

	if method == "GET" {
		query := toQuery(params)
		if query != "" {
			path += "?" + query
		}
		req, _ = http.NewRequest(method, path, nil)
	} else if len(mapx.Keys(params)) > 0 {
		if f != nil {
			pwd, _ := os.Getwd()
			file, err := os.Open(pwd + "/" + f.(string))
			if err != nil {
				panic(err)
			}
			stat, _ := file.Stat()
			buffer := &bytes.Buffer{}
			writer := multipart.NewWriter(buffer)

			for k, v := range params {
				_ = writer.WriteField(k, fmt.Sprintf("%v", v))
			}
			_, _ = writer.CreateFormFile("file", stat.Name())

			// need to know the boundary to properly close the part myself.
			boundary := writer.Boundary()
			//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
			closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

			requestReader := io.MultiReader(buffer, file, closeBuf)
			req, _ = http.NewRequest(method, path, requestReader)
			req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

			req.ContentLength = stat.Size() + int64(buffer.Len()) + int64(closeBuf.Len())
			// _ = file.Close()
		} else {
			bs, _ := json.Marshal(params)
			req, _ = http.NewRequest(method, path, bytes.NewReader(bs))
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, _ = http.NewRequest(method, path, bytes.NewReader([]byte{}))
	}

	if headers != nil {
		for k, v := range headers.(map[string]string) {
			req.Header.Set(k, v)
		}
	}
	return req
}

type Route struct {
	Path   string
	Method string
}

func GetRoute(method, path string) (route Route, found bool) {
	if path != "" && method != "" {
		return Route{Path: path, Method: method}, true
	} else if (path == "" || method == "") && CurrentAPI != "" {
		if buboEngine, ok := CurrentApp.Engine.(*bubo.App); ok {
			handlerInfo := buboEngine.GetRouterInfo().HandlerInfo[CurrentAPI]
			if handlerInfo == nil {
				return route, false
			}
			path = handlerInfo.Path
			path = strings.ReplaceAll(path, "//", "/")
			path = strings.ReplaceAll(path, "{", ":")
			path = strings.ReplaceAll(path, "}", "")
			route.Path = path
			route.Method = handlerInfo.Method
		} else {
			routeInfo := routex.NameToRouteMap[CurrentAPI]
			if routeInfo.Path == "" {
				return route, false
			}
			route.Path = routeInfo.GroupPath + routeInfo.Path
			route.Method = string(routeInfo.Method)
		}

		if method != "" {
			route.Method = method
		}
		return route, true
	} else {
		ginkgo.Fail("cannot get route automatically")
		return route, false
	}
}
