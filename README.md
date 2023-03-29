# TestX

BDD Infrastructure for Golang / Golang BDD 基础设施，关注代码可靠性与成本间的平衡

Maintainers: @will.huang  
代码覆盖率：
状态：可用，正在规划下一步演化

## Table of Contents

- [h.1 Command Line](#h1-command-line)
- [h.2 Ginkgo & Gomega 使用指南](#h2-ginkgo-gomega-使用指南)
    - [h.2.1 Ginkgo](#h21-ginkgo)
    - [h.2.2 Gomega](#h22-gomega)
- [h.3 TestX & Best Practices](#h3-testx-best-practices)
    - [h.3.1 Booting App](#h31-booting-app)
    - [h.3.2 Cleaner](#h32-cleaner)
    - [h.3.3 额外封装的 Matchers](#h33-额外封装的-matchers)
    - [h.3.4 `AssertionX` & `IsExpected()`](#h34-assertionx-isexpected)
    - [h.3.5 MonkeyPatches & Stub & Mock](#h35-monkeypatches-stub-mock)
    - [h.3.6 普通测试](#h36-普通测试)
    - [h.3.7 API 测试](#h37-api-测试)
        - [h.3.7.1 `API` block](#h371-api-block)
        - [h.3.7.2 HTTPRequest 以及对其断言](#h372-httprequest-以及对其断言)
        - [h.3.7.3 Response Matchers](#h373-response-matchers)
        - [h.3.7.4 技巧](#h374-技巧)
    - [h.3.8 Model 测试](#h38-model-测试)

## TODOs

- [ ] Doc: 工具架构与可靠性关系

## 1. Abstract

统一使用 BDD Testing 框架 [Ginkgo](https://github.com/onsi/ginkgo) 及其 Matcher [Gomega](https://github.com/onsi/gomega)

Method Stub 使用 [GoMonkey](https://github.com/agiledragon/gomonkey)，[教程](https://www.jianshu.com/p/633b55d73ddd)
测试书写遵循规范 [TODO BDD Style Guide](development_norm.md)
（临时参考：[BDD Style Guide](https://github.com/velesin/bdd-style-guide) & [Ruby Better Spec](http://www.betterspecs.org)）

WHY
1. TODO

### h.1 Command Line

请使用：（可以设置 zsh alias）
```bash
$ go test -cover -gcflags=all=-l --failFast --slowSpecThreshold=2
# or more recommended:
$ ginkgo -cover -gcflags=all=-l --failFast --slowSpecThreshold=2

# 查看覆盖情况（注：$(basename "$PWD") 是 cover 结果文件的文件名）
$ go tool cover -html=$(basename "$PWD").coverprofile
```

可留意的选项：
1. `--v`
2. `--progress`: 错误时打印 block 执行栈
3. `--focus`: 仅执行匹配到描述的测试
4. CI 有关选项：`-r` 递归执行测试 / `-outputdir=/artifacts/ -coverprofile=coverage.out` 触发 coverprofile combine

[for more information](http://onsi.github.io/ginkgo/#running-tests)

### h.2 Ginkgo & Gomega 使用指南

（此部分是官方文档的精简搬运）

#### h.2.1 要注意的问题

1. 注意唯一索引和软删除
2. 注意释放猴子补丁、还原全局变量以防止污染其他包的测试执行
3. time Format -> String 可能遇到精度不一致无法 Equal 问题
4. CI: `go get github.com/onsi/ginkgo/ginkgo@d90e0dc &&
        GIN_MODE=release ginkgo -cover -gcflags=all=-l --failFast --slowSpecThreshold=2 -r -outputdir=/artifacts/ -coverprofile=coverage.out &&
        go tool cover -func=/artifacts/coverage.out`

#### h.2.2 模板

通用的模板：（后续考虑写成 generator）
 
*_suite_test 模板
```go
import (
    "testing"

    . "github.com/go-web-kits/testx"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestXXXX(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "XXXX Suite")
}

type User struct {
    model.Default
}

var models = []interface{}{&User{}}

var _ = BeforeSuite(func() {
    BootApp(Without{Workers: true}).Migrate(models...)
})

var _ = AfterSuite(func() {
    ShutApp()
})
```

*_test 模板
```go
import (
    . "github.com/go-web-kits/testx"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("XXX", func() {
    var (
    	user User
        p *MonkeyPatches
    )

    BeforeEach(func() {
        factory.Create(&user)
    })

    AfterEach(func() {
        CleanData(&User{})
        Reset(&user)
        p.Check()
    })
	
    Describe("Func", func() {
    })
})
```

#### h.2.1 Ginkgo
#### h.2.2 Gomega

### h.3 TestX & Best Practices

TestX 提供了一系列测试辅助方法以及 BDD 封装

#### h.3.1 Booting App

```go
var _ = BeforeSuite(func() {
	// boot 配置中心、pg、redis、workers、消息网关
    BootApp()
	// boot 配置中心、pg、redis
    BootApp(Without{workers: true})
	// boot 配置中心、redis
    BootApp(Without{workers: true, Pg: true})
    
    BootApp().Migrate(&User{})
    // 同时以应用的 routes 配置启动一个测试的 Gin Server
    BootApp().Migrate(&User{}).BootGin()
})

var _ = AfterSuite(func() {
    ShutApp()
    ShutApp().Drop(&User{})
})
```

#### h.3.2 Cleaner

1. Data Cleaning
    ```go
    AfterEach(func() {
        CleanData(&User{})
    })
    ```
2. Redis Cleaning // TODO
3. Variable Cleaning: 将变量设置回零值
    ```go
    AfterEach(func() {
        Reset(&user)
    })
    ```

#### h.3.3 额外封装的 Matchers

1. `BeLike`
    ```go
    // 等效于 `BeEquivalentTo`
    Expect(1.0).To(BeLike(1))
    // slice 忽视顺序
    Expect([]interface{}{1, 2, 3.0}).To(BeLike([]int{1, 3, 2}))
    // map 仅比较给定 key-value
    Expect(map[string]interface{}{"foo": 1.0, "bar": 2}).To(BeLike(map[string]interface{}{"foo": 1}))
    // 可以各种组合
    Expect([]map[string]interface{}{{"foo": 1, "bar": 2}, {"x": "a", "y": "b"}}).
	    To(BeLike[]map[string]interface{}{{"y": "b"}, {"bar": 2}})
	// 注意 nil 返回相等
	Expect(error(nil)).To(BeLike(nil))
    ```
2. 适用于 map & struct (/ model)
    - `Include` (别名 `HaveAttributes`)
        ```go
        type H = map[string]interface{}
        Expect(H{"a": 1, "b": 2}).to(Include(H{"a": 1})) // OK
        // nested include is supported
        Expect(H{"a": H{"b": 1, "d": 2}, "x": "y"}).to(Include(H{"a": H{"b": 1}})) // OK
        
        type User struct { Name string `json:"name" db:"name"` }
        Expect(user).To(HaveAttributes(User{Name: "name"}))
        Expect(user).To(HaveAttributes(H{"name": "name"}))
        ```
        具体实现：`json.Unmarshal` 后，使用 `reflect.DeepEqual` 进行判断  
        注意：如果给定的期望值为 struct，会跳过期望 struct 的零值字段，但注意无法跳过类似 `created_at` 的字段。
        因此默认对 created_at & updated_at 做了特殊处理（跳过），其余希望跳过的字段，要主动 ignore：
        ```go
        // ignore name and id comparison
        Expect(user).To(HaveAttributes(User{Name: "name"}, "name", "id"))
        ```
3. 适用于 model instance
    - `BeTheSameRecordTo`: 比较主键是否一致
        ```go
        Expect(user).To(BeTheSameRecordTo(user1))
        ```
4. 适用于 dbx.Result
    - `HaveAffected`: 判断 dbx.Result 的 Err 是否为空
        ```go
        Expect(dbx.Update(&user)).To(HaveAffected())
        ```
    - `HaveFound`: 判断 dbx.Result 是否有 not found error
        ```go
        Expect(dbx.FindBy(&user, should.EQ{"id": 1})).To(HaveFound())
        ```
5. 适用于判断接口 response（见下文 [API 测试](#h373-response-matchers)）

#### h.3.4 `AssertionX` & `IsExpected()`

TestX 提供一种断言链，即以 AssertionX 作为接收者和返回。  
以下方法开启断言链：
- `Expectx(...)`
- `IsExpected`
- `ExpectRequested` & `ExpectRequestedBy(...)`

断言链可以做以下事情：
```go
ExpectRequested().ResponseCode().To(Equal(http.StatusOK))
ExpectRequested().ResponseBody().To(BeLike(H{"result": H{"code": 0, "message": "success"}}))
ExpectRequested().ResponseData().To(BeLike(H{"id": 1}}))
```

`IsExpected` 方法以当前测试的主体作为 Expect 的参数，并返回断言链，详细来说：
- 如果 `testx.Subject` 不为空，则以其作为 Expect 参数
- 否则如果 `testx.CurrentAPI` 不为空，则发起 `Request`
```go
// A
BeforeEach(func() {
	Subject = true
})

It("does ok", func() {
	IsExpected().To(BeTrue())
})

// B
BeforeEach(func() {
	CurrentAPI = utils.GetFuncName(user.GetHandler)
})

It("does ok", func() {
	IsExpected().To(ResponseSuccess())
	params := map[string]interface{}{"id": 1}
	IsExpected(params).To(ResponseSuccess())
})
```

#### h.3.5 MonkeyPatches & Stub & Mock

##### Simple

TestX 封装了 `gomonkey`，用反射实现运行时打猴子补丁，利用其可在运行时替换函数实现（Stub）。

使用示例如下：
1. Stub function: `IsExpectedToCall`
    ```go
    IsExpectedToCall(fmt.Sprintf).AndReturn("abc")
    // Or
    IsExpectedToCall(fmt.Sprintf).AndPerform(
        func(_ string, _ ...interface{}) string {
            // ...
        })
    ```
1. Stub method: `ExpectAnyInstanceLike`
    ```go
    ExpectAnyInstanceLike(&http.Client{}).ToCall("Do").AndReturn(
 	    &http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("body")))}, nil)
    // Or
    ExpectAnyInstanceLike(&http.Client{}).ToCall("Do").AndPerform(
        func(_ *http.Client, _ *http.Request) (*http.Response, error) {
            return &http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("body")))}, nil
        })
    ```
    注意：
    - 使用方法 stub 需要在跑 test 时增加参数 `-gcflags=all=-l` 关闭内联优化
    - 使用 `AndPerform` 传递 Stub 方法体时，作为 stub 值的匿名函数的参数列表中，第一个参数【必须】是「接收者」

提示：没有通过 `p.Reset()` / `p.Check()` 移除掉补丁，会对后续执行的测试产生影响

##### Advance

对函数（方法）调用次数进行断言：
1. 调用断言方法
2. 调用 `p.Check()`

```go
var _ = Describe("XXX", func() {
	var (
		p *MonkeyPatches
	)

    AfterEach(func() {
        p.Check()
    })

    It("does something", func() {
    	p = IsExpectedToCall(fmt.Sprintf).AndReturn("").Times(2)
    	fmt.Sprintf("")
    }) // Fail, because it only calls fmt.Sprintf Once
})
```

表示次数断言的方法有：
- `Times`
- `AtLeastOnce`
- `Once`
- `NotOnce`

注意：
1. `p = ` 不能少
2. 如果有多个 Stub，需要如此：
    ```
    p = IsExpectedToCall...
    p.IsExpectedToCall...
    p.ExpectAnyInstanceLike...
    ```

##### Mock

借由猴子补丁的能力，可以通过 Stub 实现 Mock。  
对于使用 `dbx` 的应用来说，可以直接对 `dbx` 诸方法进行 Stub，例如：  
```go
IsExpectedToCall(dbx.Where).AndReturn(dbx.Result{Data: &user})
```

`testx/let` 提供了一些快捷的 `dbx` Stub:
```go
let.UpdateBy().Succeed()
let.UpdateBy().Fail()
```

### h.3.6 普通测试

[示例](../../golang/cache/test/entry.go)

### h.3.7 API 测试

前提：使用 `routex` 进行路由注册  
（不过，不用 `routex` 也可以使用下文中部分特性）

#### h.3.7.1 `API` block

是 `Describe` 的封装，并且表示 API 测试的语义。

```go
API(HealthHandler, func() {
    It("responses successfully", func() {
    	//
    })
})
// Or (not recommended)
API("controller.HealthHandler", func() {
    It("responses successfully", func() {
    	//
    })
})
```

其主要行为是将描述（handler 函数的运行时名字）设置到 `CurrentAPI` 全局变量中。

#### h.3.7.2 HTTPRequest 以及对其断言

##### Request-Response (RR)

```go
type RR struct {
	API                string
	Params             interface{}
	ResponseCode       int
	ResponseBody       map[string]interface{}
	ResponseBodySlice  []interface{}
	ResponseBodyString string
	ResponseHeader     http.Header
}
```

##### HTTPRequest

如果没有使用 `routex` 或者想做自定义的请求，可以使用 `HTTPRequest(method, path string, param ...map[string]interface{}) RR`，  
或者它的快捷方式：`HTTPGet` / `HTTPPost` / `HTTPPut` / `HTTPDelete`
```go
r := HTTPGet("/health")
r.ResponseCode // => 200

// Query / Body 以及 Restful 参数均在同一个 map 中给定
HTTPPost("/users/:id", map[string]interface{}{"id": 1, "name": "abc"})
```

##### 如果设置了 `CurrentAPI`（即使用 `API` block）

那么可以直接使用这两个方法，其会根据 `CurrentAPI` 到路由列表中查找 path & method，进行 `HTTPRequest`
```go
CurrentAPI = "controller.HealthHandler"
r := Request()
r.ResponseCode // => 200
```
```go
API(user.Create, func() {
	It("does well", func() {
		RequestBy(H{"id": 1, "name": "abc"})
		// ...
	})
})
```

`Request` 实际上会使用 `testx.CurrentParams` 全局变量进行请求，因此你可以：
```go
CurrentParams = H{ "signature": "right", "key": "value"}
Request()
// 等同于
RequestBy(H{ "signature": "right", "key": "value"})

RequestWith(H{ "signature": "wrong" })
// 等同于
RequestBy(H{ "signature": "wrong", "key": "value"})
```

##### 发起断言

你可以使用以下三个方法发起断言：
1. `ExpectRequested`
2. `ExpectRequestedBy`
3. `ExpectRequestedWith`

```go
API(user.Create, func() {
	It("does well", func() {
		ExpectRequestedBy(H{"id": 1, "name": "abc"}).To(ResponseSuccess())
	})
})
```

#### h.3.7.3 Response Matchers

TestX 封装了一些专门用于 Request-Response 的 matcher

1. `ResponseSuccess`: 判断是否返回了 `{ "result": { "code": 0 } }` // TODO: 可配置
    ```go
    ExpectRequested().To(ResponseSuccess())
    ```
2. `Response`: 可以判断是否 response HTTPCode / Body / business_error
    ```go
    ExpectRequested().To(ResponseSuccess())
    // 等价于
    ExpectRequested().To(Response(H{"result": H{"code": 0, "message": "success"}}))
    
    // 判断 HTTP status
    ExpectRequested().To(Response(http.StatusOK))
    // 判断是否与给定 business_error 错误码相同
    ExpectRequested().To(Response(business_error.CommonError[business_error.NotFound]))
    // 判断是否与给定 error 有相同 message
    ExpectRequested().To(Response(errors.New("error msg")))

    // 判断 Body map 是否 `BeLike`
    ExpectRequested().To(Response(H{"data": H{"id": 1, "name": "abc"}}))
    // 或者
    ExpectRequested().To(ResponseData(H{"id": 1, "name": "abc"}))
    ```

#### h.3.7.4 技巧

1. 以自定义的 routes 配置启动 Gin，并能够在每个 test spec 内动态修改 handler 行为
    ```go
    // package response_test
    
    var Action func(*gin.Context)
    var Handler = func(c *gin.Context) { Action(c) }
    
    var _ = BeforeSuite(func() {
        IsExpectedToCall(routes.InitRoutes).AndPerform(func() {
            routes.Routes = []interface{}{
                routex.GETx("/", Handler),
            }
        })
        BootApp().BootGin()
    })
    
    var _ = Describe("Success", func() {
        BeforeEach(func() {
            CurrentAPI = utils.GetFuncName(Handler)
        })
    
        It("responses success", func() {
            Action = func(c *gin.Context) {
                // Your Action
            }
            IsExpected().To(ResponseSuccess())
        })
    })
    ```

### h.3.8 Model 测试

factory
