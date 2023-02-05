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
	authRouter.GET("/get-room-request-warden", func(c *gin.Context) {
		api.GetAllRoomRequests(c, &database)
	})
	authRouter.GET("/get-requested-room", func(c *gin.Context) {
		api.GetRequestedRoom(c, &database)
	})
	authRouter.POST("/request-room", func(c *gin.Context) {
		api.RequestRoom(c, &database)
	})
	authRouter.POST("/accept-room", func(c *gin.Context) {
		api.AcceptRoom(c, &database)
	})
	authRouter.POST("/accept-room-warden", func(c *gin.Context) {
		api.AcceptRoomWarden(c, &database)
	})

	router.Run(":8000")

}
