package testx

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/spear/bubo"
	"github.com/onsi/ginkgo"
)

type Engine interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type App struct {
	Engine
}

var (
	CurrentApp *App
)

func BootApp(f func(args ...interface{}) error, args ...interface{}) *App {
	_ = os.Setenv("TEST", "true")

	if f(args...) != nil {
		ginkgo.Fail("Cannot Boot App")
	}
	CurrentApp = &App{}
	return CurrentApp
}

func Boot() *App {
	_ = os.Setenv("TEST", "true")
	CurrentApp = &App{}
	return CurrentApp
}

func ShutApp() *App {
	_ = os.Setenv("TEST", "false")
	return CurrentApp
}

func (app *App) Migrate(values ...interface{}) *App {
	// dbx.Conn(dbx.Opt{UnLog: false}).DropTableIfExists(values...)
	dbx.Conn(dbx.Opt{UnLog: false}).AutoMigrate(values...)
	// if err != nil {
	// 	panic(err)
	// }
	return app
}

func (app *App) Drop(values ...interface{}) *App {
	dbx.Conn(dbx.Opt{UnLog: true}).DropTableIfExists(values...)
	return app
}

func (app *App) BootGin(e *gin.Engine) *App {
	app.Engine = e
	return app
}

func (app *App) BootSpear(e *bubo.App) *App {
	app.Engine = e
	return app
}
