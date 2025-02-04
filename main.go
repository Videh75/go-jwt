package main

import (
	"go-jwt/controllers"
	"go-jwt/database"
	"go-jwt/initializers"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	database.ConnectToDB()
	r := gin.Default()
	r.POST("/signup", func(c *gin.Context) {
		controllers.SignUp(c)
	})
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c)
	})
	r.GET("/validate", func(c *gin.Context) {
		middleware.RequireAuth(c)
	}, func(c *gin.Context) {
		controllers.Validate(c)
	})
	r.Run()
}
