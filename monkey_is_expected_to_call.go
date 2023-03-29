package testx

import (
	"reflect"

	"github.com/agiledragon/gomonkey"
)

type Chain1 struct {
	target interface{}
	monkey *MonkeyPatches
}

func IsExpectedToCall(target interface{}) Chain1 {
	return Chain1{target: target}
}

func (c Chain1) AndPerform(double interface{}) (monkey *MonkeyPatches) {
	targetType := reflect.TypeOf(c.target)
	if c.monkey != nil {
		monkey = c.monkey.addMonkeyFor(c.target)
		monkey.Patches = monkey.ApplyFunc(c.target, monkey.WrapDouble(targetType, double))
	} else {
		monkey = newMonkeyFor(c.target)
		monkey.Patches = gomonkey.ApplyFunc(c.target, monkey.WrapDouble(targetType, double))
	}
	return
}

func (c Chain1) AndReturn(vals ...interface{}) (monkey *MonkeyPatches) {
	targetType := reflect.TypeOf(c.target)
	if c.monkey != nil {
		monkey = c.monkey.addMonkeyFor(c.target)
		monkey.Patches = monkey.ApplyFunc(c.target, monkey.MakeDouble(targetType, vals...))
	} else {
		monkey = newMonkeyFor(c.target)
		monkey.Patches = gomonkey.ApplyFunc(c.target, monkey.MakeDouble(targetType, vals...))
	}
	return
}

func (c Chain1) JustReturn(vals ...interface{}) (monkey *MonkeyPatches) {
	targetType := reflect.TypeOf(c.target)
	if c.monkey != nil {
		monkey = c.monkey.addMonkeyFor(c.target)
		monkey.Patches = monkey.ApplyFunc(c.target, monkey.JustMakeDouble(targetType, vals...))
	} else {
		monkey = newMonkeyFor(c.target)
		monkey.Patches = gomonkey.ApplyFunc(c.target, monkey.JustMakeDouble(targetType, vals...))
	}
	return
}
