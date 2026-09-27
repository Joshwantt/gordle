package session

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

const maxGuesses = 6

func CreateSession(sessionStore *store.Store) gin.HandlerFunc {
	return func(context *gin.Context) {
		hiddenWord := "crate"
		newSession := gamesession.New(hiddenWord, maxGuesses)

		sessionID := sessionStore.AddSession(newSession)

		context.JSON(http.StatusCreated, NewResponse(sessionID, newSession))
	}
}
