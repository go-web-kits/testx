package testx

import (
	"reflect"

	"github.com/go-web-kits/dbx"
)

// TODO
// https://jarifibrahim.github.io/blog/test-cleanup-with-gorm-hooks/
func CleanData(values ...interface{}) {
	// ginkgo.BeforeEach(func() {
	// 	for _, v := range values {
	// 		dbx.Destroy(v, dbx.Opt{SkipCallback: true, Unscoped: true})
	// 	}
	// })

	// ginkgo.AfterEach(func() {
	for _, v := range values {
		dbx.DestroyAll(v, dbx.Opt{SkipCallback: true, Unscoped: true})
	}
	// })
}

func Reset(values ...interface{}) {
	for _, value := range values {
		v := reflect.ValueOf(value).Elem()
		v.Set(reflect.Zero(v.Type()))
	}
}
