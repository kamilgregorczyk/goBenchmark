package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

type Post struct {
	Title     string
	CreatedAt string
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("error")
	app.RegisterView(iris.HTML("./templates", ".html"))
	app.Get("/", func(ctx context.Context) {
		ctx.ViewData("Title", "Nice title")
		ctx.ViewData("Posts", []Post{
			{Title: "ABC", CreatedAt: "2018-01-01"},
			{Title: "DEF", CreatedAt: "2018-01-01"},
			{Title: "GHI", CreatedAt: "2018-01-01"},
		})
		ctx.View("index.html")
	})
	app.Run(iris.Addr(":8000"))
}
