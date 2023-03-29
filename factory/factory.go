package factory

import (
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/utils"
	"github.com/go-web-kits/utils/slicex"
	"github.com/jinzhu/gorm"
	"github.com/onsi/ginkgo"
)

func Create(obj interface{}, opts ...dbx.Opt) interface{} {
	slicex.Each(obj, func(item interface{}) {
		assert("create", item, dbx.Create(item, opts...))
	})
	return obj
}

func UpdateBy(obj interface{}, values interface{}, opts ...dbx.Opt) interface{} {
	slicex.Each(obj, func(item interface{}) {
		assert("update", item, dbx.UpdateBy(item, values, opts...))
	})
	return obj
}

func Destroy(obj interface{}, opts ...dbx.Opt) interface{} {
	slicex.Each(obj, func(item interface{}) {
		assert("destroy", item, dbx.Destroy(item, opts...))
	})
	return obj
}

func Association(obj interface{}, column string) *gorm.Association {
	return dbx.Conn().Model(obj).Association(column)
}

func Find(obj interface{}, opts ...dbx.Opt) interface{} {
	assert("find", obj, dbx.Find(obj, nil, opts...))
	return obj
}

func Reload(obj interface{}, opts ...dbx.Opt) {
	slicex.Each(obj, func(item interface{}) {
		assert("reload", item, dbx.FindById(item, dbx.IdOf(item), opts...))
	})
}

// =====

func assert(action string, obj interface{}, result dbx.Result) dbx.Result {
	if result.Err != nil {
		ginkgo.Fail(utils.TypeNameOf(obj) + " " + action + " failed. error: " + result.Err.Error())
	}
	return result
}
