package testx

import "github.com/agiledragon/gomonkey"

type Chain3 struct {
	target  interface{}
	patches *gomonkey.Patches
}

func ExpectGlobalVar(value interface{}) Chain3 {
	return Chain3{target: value}
}

func (c Chain3) ToBe(double interface{}) *MonkeyPatches {
	if c.patches != nil {
		return &MonkeyPatches{Patches: c.patches.ApplyGlobalVar(c.target, double)}
	}
	return &MonkeyPatches{Patches: gomonkey.ApplyGlobalVar(c.target, double)}
}
