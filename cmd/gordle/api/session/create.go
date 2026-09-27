package session

import (
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"

	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

// hard coded 6 guesses for session construction.
const maxGuesses = 6

// hard coded 5 letter words for moment
const wordLength = 5

// start a new session with a random wordLength hidden word from answersByLength and respond with its initial state
func CreateSession(sessionStore *store.Store, answersByLength map[int][]string) gin.HandlerFunc {
	return func(context *gin.Context) {
		answers := answersByLength[wordLength]
		hiddenWord := answers[rand.N(len(answers))]
		newSession := gamesession.New(hiddenWord, maxGuesses)

		sessionID := sessionStore.AddSession(newSession)

		context.JSON(http.StatusCreated, NewResponse(sessionID, newSession))
	}
}
