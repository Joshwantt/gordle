// Package api implements Gordle's HTTP API.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api/health"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session/session_id"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

// the largest request body any route accepts.
const maxRequestBodyBytes = 2048

// attach every API route to routerGroup, with session routes backed by sessionStore
func Register(routerGroup *gin.RouterGroup, sessionStore *store.Store, answersByLength map[int][]string, guessesByLength map[int]map[string]bool) {
	routerGroup.Use(func(context *gin.Context) {
		context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, maxRequestBodyBytes)
		context.Next()
	})

	routerGroup.GET("/health", health.Get)

	sessionGroup := routerGroup.Group("/session")
	sessionGroup.POST("", session.CreateSession(sessionStore, answersByLength))
	sessionGroup.GET("/:session_id", session_id.GetSession(sessionStore))
	sessionGroup.POST("/:session_id/guess", session_id.MakeGuess(sessionStore, guessesByLength))
}
