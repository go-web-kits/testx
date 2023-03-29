package testx

import (
	"reflect"

	"github.com/agiledragon/gomonkey"
)

type Chain2 struct {
	targetType reflect.Type
	method     string
	monkey     *MonkeyPatches
}

func ExpectAnyInstanceLike(target interface{}) Chain2 {
	return Chain2{targetType: reflect.TypeOf(target)}
}

func (c Chain2) ToCall(method string) Chain2 {
	c.method = method
	return c
}

func (c Chain2) AndPerform(double interface{}) (monkey *MonkeyPatches) {
	method, ok := c.targetType.MethodByName(c.method)
	if !ok {
		panic("no such method that you wanna to stub")
	}
	if c.monkey != nil {
		monkey = c.monkey.addMonkeyFor(c.method)
		monkey.Patches = monkey.ApplyMethod(c.targetType, c.method, monkey.WrapDouble(method.Type, double))
	} else {
		monkey = newMonkeyFor(c.method)
		monkey.Patches = gomonkey.ApplyMethod(c.targetType, c.method, monkey.WrapDouble(method.Type, double))
	}
	return
}

func (c Chain2) AndReturn(vals ...interface{}) (monkey *MonkeyPatches) {
	method, ok := c.targetType.MethodByName(c.method)
	if !ok {
		panic("no such method that you wanna to stub")
	}
	if c.monkey != nil {
		monkey = c.monkey.addMonkeyFor(c.method)
		monkey.Patches = monkey.ApplyMethod(c.targetType, c.method, monkey.MakeDouble(method.Type, vals...))
	} else {
		monkey = newMonkeyFor(c.method)
		monkey.Patches = gomonkey.ApplyMethod(c.targetType, c.method, monkey.MakeDouble(method.Type, vals...))
	}
	return
}
