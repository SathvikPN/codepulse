package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomeHandler(c *gin.Context) {
	// err := InsertUser(c.Query("name"), c.ClientIP())
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "internal error at server")
	// }
	resp := map[string]interface{}{
		"remoteAddr": c.Request.RemoteAddr,
		"remoteIP":   c.RemoteIP(),
		"uri":        c.Request.URL,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
