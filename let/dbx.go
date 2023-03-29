package let

import (
	"errors"

	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/testx"
)

type DBx struct {
	Target interface{}
}

func (x DBx) Succeed() *testx.MonkeyPatches {
	return testx.IsExpectedToCall(x.Target).AndReturn(dbx.Result{})
}

func (x DBx) Fail() *testx.MonkeyPatches {
	return testx.IsExpectedToCall(x.Target).AndReturn(dbx.Result{Err: errors.New("")})
}

func UpdateBy() DBx {
	return DBx{Target: dbx.UpdateBy}
}
