package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

func main() {
	app := iris.New()

	app.StaticWeb("/static", "./resources")

	app.RegisterView(iris.HTML("./templates", ".html").Layout("layout.html").Reload(true))

	app.Get("/", func(ctx context.Context) {
		ctx.Gzip(true)
		ctx.View("cla.html")
	})

	app.Run(iris.Addr(":8080"))
}
