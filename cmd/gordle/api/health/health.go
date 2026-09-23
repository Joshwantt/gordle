// Package health implements the /api/health route.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// are we alive?
func Get(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "ok"})
}
