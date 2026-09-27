// Package api implements Gordle's HTTP API.
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api/health"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session/session_id"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

// attach every API route to routerGroup, with session routes backed by sessionStore
func Register(routerGroup *gin.RouterGroup, sessionStore *store.Store, answersByLength map[int][]string, guessesByLength map[int]map[string]bool) {
	routerGroup.GET("/health", health.Get)

	sessionGroup := routerGroup.Group("/session")
	sessionGroup.POST("", session.CreateSession(sessionStore, answersByLength))
	sessionGroup.GET("/:session_id", session_id.GetSession(sessionStore))
	sessionGroup.POST("/:session_id/guess", session_id.MakeGuess(sessionStore, guessesByLength))
}
