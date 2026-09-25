package session_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// make a guess in an existing session
func MakeGuess(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{})
}
