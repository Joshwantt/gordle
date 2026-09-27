package session_id

import (
	"net/http"
	"strings"
	"uuid"

	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
	"github.com/gin-gonic/gin"
)

type guessRequest struct {
	Guess string `json:"guess" binding:"required"`
}

// make a guess in an existing session accepting only words from guessesByLength
func MakeGuess(sessionStore *store.Store, guessesByLength map[int]map[string]bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionID, parseError := uuid.Parse(context.Param("session_id"))

		if parseError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"reason": "invalid_request"})
			return
		}

		var request guessRequest

		bindError := context.ShouldBindJSON(&request)
		if bindError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"reason": "invalid_request"})
			return
		}

		guessWord := strings.ToLower(request.Guess)

		// UpdateSession will hold the mutex
		// Capture outcomes so they can be replied to after UpdateSession has exited (and released the mutex)
		var failureStatus int
		var failureReason string
		var updatedSession gamesession.Session

		exists := sessionStore.UpdateSession(sessionID, func(storedSession *gamesession.Session) {
			if storedSession.IsComplete() {
				failureStatus, failureReason = http.StatusConflict, "session_complete"
				return
			}

			if !guessesByLength[storedSession.WordLength][guessWord] {
				failureStatus, failureReason = http.StatusUnprocessableEntity, "invalid_word"
				return
			}

			storedSession.Guesses = append(storedSession.Guesses, gamesession.Guess{
				Word:   guessWord,
				Result: gamesession.Score(guessWord, storedSession.HiddenWord),
			})
			updatedSession = *storedSession
		})

		// UpdateSession returns a bool which is only false if the session is not found
		if !exists {
			context.JSON(http.StatusNotFound, gin.H{"reason": "session_not_found"})
			return
		}

		if failureStatus != 0 {
			context.JSON(failureStatus, gin.H{"reason": failureReason})
			return
		}

		context.JSON(http.StatusOK, session.NewResponse(sessionID, updatedSession))
	}
}
