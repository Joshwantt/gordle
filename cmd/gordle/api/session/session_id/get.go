package session_id

import (
	"net/http"
	"uuid"

	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
	"github.com/gin-gonic/gin"
)

// look up an existing session by its session_id
func GetSession(sessionStore *store.Store) gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionID, parseError := uuid.Parse(context.Param("session_id"))
		if parseError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"reason": "invalid_request"})
			return
		}

		storedSession, exists := sessionStore.GetSession(sessionID)
		if !exists {
			context.JSON(http.StatusNotFound, gin.H{"reason": "session_not_found"})
			return
		}

		context.JSON(http.StatusOK, session.NewResponse(sessionID, storedSession))
	}
}
