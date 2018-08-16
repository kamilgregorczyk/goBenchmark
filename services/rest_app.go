package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func main() {
	app := iris.Default()
	app.Logger().SetLevel("error")
	app.Use(logger.New())
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(map[string]string{"hello": "json"})
	})
	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(":8000"))
}