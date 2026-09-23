// Package api implements Gordle's HTTP API.
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api/health"
)

func Register(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/health", health.Get)
}
