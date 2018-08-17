package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gopkg.in/go-playground/validator.v9"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"context"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

// Models

type Post struct {
	ID        objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string            `bson:"title" json:"title" validate:"required"`
	CreatedAt string            `bson:"created_at" json:"created_at"`
}

type User struct {
	ID       objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string            `json:"email" validate:"required"`
	Password string            `json:"password" validate:"required"`
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Controllers

type UsersController struct {
	Ctx             iris.Context
	Validate        *validator.Validate
	UsersCollection *mongo.Collection
}

func (c *UsersController) Post() {
	var user User
	if err := c.Ctx.ReadJSON(&user); err != nil {
		c.Ctx.StatusCode(400)
		c.Ctx.WriteString(err.Error())
		return
	}
	err := c.Validate.Struct(user)
	if err != nil {
		c.Ctx.StatusCode(400)
		c.Ctx.WriteString(err.Error())
		return
	}

	count, err := c.UsersCollection.Count(nil, bson.NewDocument(bson.EC.String("username", user.Email)))
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.WriteString(err.Error())
		return
	}
	if count >= 1 {
		c.Ctx.StatusCode(400)
		c.Ctx.WriteString("User already exists!")
		return
	}
	result, err := c.UsersCollection.InsertOne(nil, bson.NewDocument(
		bson.EC.String("username", user.Email),
		bson.EC.String("password", HashPassword(user.Password)),
	))
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.StatusCode(201)
	c.Ctx.JSON(result)
}

type PostsController struct {
	Ctx             iris.Context
	Validate        *validator.Validate
	PostsCollection *mongo.Collection
}

func (c *PostsController) Get() {
	cursor, err := c.PostsCollection.Find(context.Background(), nil)
	if err != nil {
		c.Ctx.StatusCode(500)
		return
	}
	var posts []Post
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		post := Post{}
		err := cursor.Decode(&post)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)

	}
	c.Ctx.JSON(posts)

}
func (c *PostsController) Post() {
	var post Post
	if err := c.Ctx.ReadJSON(&post); err != nil {
		c.Ctx.StatusCode(400)
		c.Ctx.WriteString(err.Error())
		return
	}
	err := c.Validate.Struct(post)
	if err != nil {
		c.Ctx.StatusCode(400)
		c.Ctx.WriteString(err.Error())
		return
	}
	post.CreatedAt = time.Now().String()
	result, err := c.PostsCollection.InsertOne(nil, bson.NewDocument(
		bson.EC.String("title", post.Title),
		bson.EC.String("created_at", post.CreatedAt),
	))
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.StatusCode(201)
	c.Ctx.JSON(result)
}

func (c *PostsController) GetBy(id string) {
	post := Post{}
	object_id, err := objectid.FromHex(id)
	if err != nil {
		c.Ctx.StatusCode(400)
		return
	}
	doc := c.PostsCollection.FindOne(nil, bson.NewDocument(
		bson.EC.ObjectID("_id", object_id),
	))
	if doc == nil {
		c.Ctx.StatusCode(404)
		return
	}
	doc.Decode(&post)
	c.Ctx.JSON(post)
}

func (c *PostsController) DeleteBy(id string) {
	object_id, err := objectid.FromHex(id)
	if err != nil {
		c.Ctx.StatusCode(400)
		return
	}
	c.PostsCollection.DeleteOne(nil, bson.NewDocument(
		bson.EC.ObjectID("_id", object_id),
	))
}

var validate *validator.Validate

func main() {
	validate = validator.New()
	app := iris.New()
	app.Logger().SetLevel("error")
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	postsRouter := app.Party("/posts")
	postsRouter.Use(jwtHandler.Serve)
	posts := mvc.New(postsRouter)
	posts.Register(client.Database("go").Collection("posts"))
	posts.Register(validate)
	posts.Handle(new(PostsController))

	users := mvc.New(app.Party("/users"))
	users.Register(client.Database("go").Collection("users"))
	users.Register(validate)
	users.Handle(new(UsersController))

	app.Run(iris.Addr(":8000"),
		iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
