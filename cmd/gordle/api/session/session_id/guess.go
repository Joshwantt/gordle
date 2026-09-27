package session_id

import (
	"net/http"
	"slices"
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
func MakeGuess(sessionStore *store.Store, guessesByLength map[int][]string) gin.HandlerFunc {
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

		exists := sessionStore.UpdateSession(sessionID, func(storedSession *gamesession.Session) {
			if storedSession.IsComplete() {
				context.JSON(http.StatusConflict, gin.H{"reason": "session_complete"})
				return
			}

			if !slices.Contains(guessesByLength[storedSession.WordLength], guessWord) {
				context.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "invalid_word"})
				return
			}

			storedSession.Guesses = append(storedSession.Guesses, gamesession.Guess{
				Word:   guessWord,
				Result: gamesession.Score(guessWord, storedSession.HiddenWord),
			})
			context.JSON(http.StatusOK, session.NewResponse(sessionID, *storedSession))
		})

		// UpdateSession returns a bool which is only false if the session is not found
		if !exists {
			context.JSON(http.StatusNotFound, gin.H{"reason": "session_not_found"})
		}
	}
}
