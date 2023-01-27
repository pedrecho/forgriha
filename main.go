package main

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"schedule/controllers"
	"schedule/middlewares"
	"schedule/models"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/static", "./static")

	public := r.Group("/")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)

	r.Run(":8080")

}
