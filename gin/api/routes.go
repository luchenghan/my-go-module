package api

import (
	"mymodule/gin/internal/handlers"
	"mymodule/gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)

	api := r.Group("/api")
	api.POST("/user", handlers.CreateUser)
	api.Use(middleware.JWTAuthorization())
	{
		api.GET("/user/:id", handlers.GetUserByID)
		api.GET("/users", handlers.GetUsers)
	}

}
