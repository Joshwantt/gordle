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

// the game is over once the hidden word is guessed or every guess is used
func (session *Session) IsComplete() bool {
	if len(session.Guesses) >= session.MaxGuesses {
		return true
	}
	return len(session.Guesses) > 0 && session.Guesses[len(session.Guesses)-1].Word == session.HiddenWord
}

// i didn't know this was how wordle's scoring function works
// i googled and adapted this one ngl
func Score(guessWord string, hiddenWord string) []LetterResult {
	results := make([]LetterResult, len(guessWord))
	unmatchedLetterCounts := map[byte]int{}

	for index := range guessWord {
		if guessWord[index] == hiddenWord[index] {
			results[index] = Correct
		} else {
			unmatchedLetterCounts[hiddenWord[index]]++
		}
	}

	for index := range guessWord {
		if results[index] == Correct {
			continue
		}
		if unmatchedLetterCounts[guessWord[index]] > 0 {
			results[index] = Present
			unmatchedLetterCounts[guessWord[index]]--
		} else {
			results[index] = Absent
		}
	}

	return results
}
