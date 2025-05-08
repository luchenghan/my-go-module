package main

import (
	"fmt"
	"mymodule/gin/api"
	"mymodule/gin/internal/config"
	"mymodule/gin/internal/db"
	"mymodule/gin/internal/services/users"
	"mymodule/gin/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger.Initialize(config.GlobalConfig.Logger.Encoding, config.GlobalConfig.Logger.Level)

	if err := db.InitMySQL(config.GlobalConfig.Database); err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	users.Init()

	// Set Gin mode from config
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// Initialize router
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router)

	// Run server with configured port
	port := fmt.Sprintf(":%s", config.GlobalConfig.Server.Port)
	if err := router.Run(port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
