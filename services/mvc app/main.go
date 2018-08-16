package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/hero"
)

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HashPassword(password string) string {
	return password
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("error")
	app.Run(iris.Addr(":8000"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
