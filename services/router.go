package services

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	api := router.Group("/api")

	api.GET("/welcome", WelcomeHandler)

	return router
}
