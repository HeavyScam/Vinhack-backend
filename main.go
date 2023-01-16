package main

import (
	// "fmt"

	"net/http"
	// "sync"
	//import the token package
	"hostel_hopper/api"
	"hostel_hopper/infrastructure"
	"hostel_hopper/token"

	// import gin
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(api.CorsMiddleware())
	infrastructure.LoadEnv() //loading env
	database := infrastructure.NewDatabase()
	tokenMaker, _ := token.InitializePasetoToken()

	authRouter := router.Group("/auth").Use(api.AuthMiddleware(tokenMaker))

	//create a new route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//register
	router.POST("/register", func(c *gin.Context) {
		api.Register(c, &database, tokenMaker)
	})

	//login
	router.POST("/login", func(c *gin.Context) {
		api.Login(c, &database, tokenMaker)
	})
	authRouter.POST("/add-preferences", func(c *gin.Context) {
		api.AddPreferences(c, &database)
	})
	authRouter.GET("/get-rooms", func(c *gin.Context) {
		api.GetAllRooms(c, &database)
	})

	router.Run(":8000")

}
