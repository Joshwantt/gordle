// Package session implements the /api/session routes.
package session

import (
	"uuid"

	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
)

// the JSON body shared by every session route
type Response struct {
	SessionID  string              `json:"session_id"`
	Guesses    []gamesession.Guess `json:"guesses"`
	MaxGuesses int                 `json:"max_guesses"`
	WordLength int                 `json:"word_length"`
	HiddenWord string              `json:"hidden_word,omitempty"`
}

// build the Response for a session, keeping the hidden word out of it
func NewResponse(sessionID uuid.UUID, gameSession gamesession.Session) Response {
	response := Response{
		SessionID:  sessionID.String(),
		Guesses:    gameSession.Guesses,
		MaxGuesses: gameSession.MaxGuesses,
		WordLength: gameSession.WordLength,
	}
	if gameSession.IsComplete() {
		response.HiddenWord = gameSession.HiddenWord
	}
	return response
}
