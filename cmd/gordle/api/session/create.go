package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// start a new session
func CreateSession(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{})
}
