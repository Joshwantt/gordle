// Package session implements the /api/session routes.
package session

import (
	"uuid"

	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
)

type Response struct {
	SessionID  string              `json:"session_id"`
	Guesses    []gamesession.Guess `json:"guesses"`
	MaxGuesses int                 `json:"max_guesses"`
	WordLength int                 `json:"word_length"`
}

func NewResponse(sessionID uuid.UUID, gameSession gamesession.Session) Response {
	return Response{
		SessionID:  sessionID.String(),
		Guesses:    gameSession.Guesses,
		MaxGuesses: gameSession.MaxGuesses,
		WordLength: gameSession.WordLength,
	}
}
