// Package session implements Gordle's game rules
package session

import "time"

type LetterResult string

const (
	Correct LetterResult = "correct"
	Present LetterResult = "present"
	Absent  LetterResult = "absent"
)

type Guess struct {
	Word   string         `json:"word"`
	Result []LetterResult `json:"result"`
}

type Session struct {
	Guesses         []Guess
	MaxGuesses      int
	LastInteraction time.Time
	HiddenWord      string
	WordLength      int
}

func New(hiddenWord string, maxGuesses int) Session {
	return Session{
		HiddenWord: hiddenWord,
		MaxGuesses: maxGuesses,
		WordLength: len(hiddenWord),
		Guesses:    []Guess{},
	}
}
