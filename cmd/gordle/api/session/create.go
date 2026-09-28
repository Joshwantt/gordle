package session

import (
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"

	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

type createRequest struct {
	WordLength int `json:"word_length" binding:"required,min=3,max=16"`
	MaxGuesses int `json:"max_guesses" binding:"required,min=1,max=20"`
}

// start a new session with a random hidden word of the requested length from answersByLength and respond with its initial state
func CreateSession(sessionStore *store.Store, answersByLength map[int][]string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var request createRequest

		bindError := context.ShouldBindJSON(&request)
		if bindError != nil {
			context.JSON(http.StatusBadRequest, gin.H{"reason": "invalid_request"})
			return
		}

		answers := answersByLength[request.WordLength]
		hiddenWord := answers[rand.N(len(answers))]
		newSession := gamesession.New(hiddenWord, request.MaxGuesses)

		sessionID := sessionStore.AddSession(newSession)

		context.JSON(http.StatusCreated, NewResponse(sessionID, newSession))
	}
}
