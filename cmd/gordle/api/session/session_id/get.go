package session_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// look up an existing session by its session_id
func GetSession(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{})
}
