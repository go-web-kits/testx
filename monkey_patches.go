package testx

import (
	"fmt"
	"reflect"

	"github.com/agiledragon/gomonkey"
	"github.com/go-web-kits/utils"
	"github.com/go-web-kits/utils/reflectx"
	"github.com/onsi/ginkgo"
)

type MonkeyPatches struct {
	*gomonkey.Patches
	CurFuncName     string
	expectFuncTimes map[string]int
	actualFuncTimes map[string]int
}

type MP = MonkeyPatches

var NotOnce = -1
var AtLeastOnce = 9999

func (p *MonkeyPatches) Check() {
	if p == nil {
		return
	}
	p.Reset()

	for fn, actualTimes := range p.actualFuncTimes {
		expectTimes := p.expectFuncTimes[fn]
		if expectTimes == 0 {
			continue
		}
		if expectTimes == NotOnce {
			if actualTimes != 0 {
				ginkgo.Fail(
					fmt.Sprintf(
						"Expect the method `%s` \n\n not to be called, \n\n but actually it had been called -> %v <- times",
						fn, actualTimes,
					),
				)
			} else {
				continue
			}
		}
		if (expectTimes == AtLeastOnce && actualTimes < 1) ||
			(expectTimes != AtLeastOnce && actualTimes != expectTimes) {
			ginkgo.Fail(
				fmt.Sprintf(
					"Expect the method `%s` \n\n to be called -> %v <- times, \n\n but actual times is -> %v <-",
					fn, expectTimes, actualTimes,
				),
			)
		}
	}

	p.expectFuncTimes = map[string]int{}
	p.actualFuncTimes = map[string]int{}
}

func (p *MonkeyPatches) IsExpectedToCall(target interface{}) Chain1 {
	return Chain1{target: target, monkey: p}
}

func (p *MonkeyPatches) ExpectAnyInstanceLike(target interface{}) Chain2 {
	return Chain2{targetType: reflect.TypeOf(target), monkey: p}
}

func (p *MonkeyPatches) NotOnce() *MonkeyPatches {
	p.expectFuncTimes[p.CurFuncName] = NotOnce
	return p
}

func (p *MonkeyPatches) Once() *MonkeyPatches {
	p.expectFuncTimes[p.CurFuncName] = 1
	return p
}

func (p *MonkeyPatches) AtLeastOnce() *MonkeyPatches {
	p.expectFuncTimes[p.CurFuncName] = AtLeastOnce
	return p
}

func (p *MonkeyPatches) Times(num int) *MonkeyPatches {
	p.expectFuncTimes[p.CurFuncName] = num
	return p
}

func (p *MonkeyPatches) On(fn string) {
	oldTimes := p.actualFuncTimes[fn]
	p.actualFuncTimes[fn] = oldTimes + 1
}

// ==========

func newMonkeyFor(target interface{}) *MonkeyPatches {
	funcName := utils.GetFuncName(target)
	eMap, aMap := map[string]int{}, map[string]int{}
	eMap[funcName] = 0
	aMap[funcName] = 0
	return &MonkeyPatches{CurFuncName: funcName, expectFuncTimes: eMap, actualFuncTimes: aMap}
}

func (p *MonkeyPatches) addMonkeyFor(target interface{}) *MonkeyPatches {
	p.CurFuncName = utils.GetFuncName(target)
	p.expectFuncTimes[p.CurFuncName] = 0
	p.actualFuncTimes[p.CurFuncName] = 0
	return p
}

func (p *MonkeyPatches) WrapDouble(targetType reflect.Type, double interface{}) interface{} {
	warpedDouble := func(in []reflect.Value) []reflect.Value {
		p.On(p.CurFuncName)
		// FIXME: variadic params
		return reflect.ValueOf(double).Call(in)
	}
	doubled := reflect.MakeFunc(targetType, warpedDouble)
	return doubled.Interface()
}

func (p *MonkeyPatches) MakeDouble(targetType reflect.Type, returns ...interface{}) interface{} {
	curFuncName := p.CurFuncName
	double := func(in []reflect.Value) []reflect.Value {
		p.On(curFuncName)

		v := []reflect.Value{}
		for i, ret := range returns {
			// FIXME: wrong type for interface{}
			v = append(v, reflectx.ValueOf(ret, targetType.Out(i)))
		}
		return v
	}

	doubled := reflect.MakeFunc(targetType, double)
	return doubled.Interface()
}

func (p *MonkeyPatches) JustMakeDouble(targetType reflect.Type, returns ...interface{}) interface{} {
	double := func(in []reflect.Value) []reflect.Value {
		v := []reflect.Value{}
		for i, ret := range returns {
			// FIXME: wrong type for interface{}
			v = append(v, reflectx.ValueOf(ret, targetType.Out(i)))
		}
		return v
	}

	doubled := reflect.MakeFunc(targetType, double)
	return doubled.Interface()
}
