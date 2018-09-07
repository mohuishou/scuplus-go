package middleware

import (
	"github.com/iris-contrib/middleware/newrelic"
	"github.com/kataras/iris"
	"github.com/mohuishou/scuplus-go/config"
)

// Register 注册中间件
func Register(app *iris.Application) {
	conf := config.Get().NewRelic
	c := newrelic.Config(conf.AppName, conf.Key)
	c.Enabled = true
	m, err := newrelic.New(c)
	if err != nil {
		app.Logger().Fatal(err)
	}
	app.Use(m.ServeHTTP)
	app.Use(jwtMiddle)
}
