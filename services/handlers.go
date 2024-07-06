package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomeHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "anonymous"
	}
	err := InsertUser(name, c.ClientIP())
	if err != nil {
		log.Printf("[ERROR] (inserting user) - %v", err)
		c.JSON(http.StatusInternalServerError, "[Internal-Error]")
		return
	}

	resp := map[string]interface{}{
		"remoteAddr": c.Request.RemoteAddr,
		"remoteIP":   c.RemoteIP(),
		"uri":        c.Request.URL,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func compareHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, nil)
}
