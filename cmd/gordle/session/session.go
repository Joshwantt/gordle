// Package session implements Gordle's game rules
package session

import "time"

type LetterResult string

const (
	// right letter in the right position. Green
	Correct LetterResult = "correct"
	// letter is in the hidden word but in a different position. Yellow
	Present LetterResult = "present"
	// letter is not in the hidden word. Grey
	Absent LetterResult = "absent"
)

type Guess struct {
	Word   string         `json:"word"`
	Result []LetterResult `json:"result"`
}

// the full state of one game, including the hidden word
type Session struct {
	Guesses         []Guess
	MaxGuesses      int
	LastInteraction time.Time
	HiddenWord      string
	WordLength      int
}

// start a fresh session
func New(hiddenWord string, maxGuesses int) Session {
	return Session{
		HiddenWord: hiddenWord,
		MaxGuesses: maxGuesses,
		WordLength: len(hiddenWord),
		Guesses:    []Guess{},
	}
}
