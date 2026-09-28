package session_id_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/Joshwantt/gordle/cmd/gordle/api/apitest"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
	gamesession "github.com/Joshwantt/gordle/cmd/gordle/session"
)

// create a session on router whose hidden word is always "crate"
func newGuessPath(router http.Handler) string {
	var created session.Response
	json.Unmarshal(apitest.ServeBody(router, http.MethodPost, "/api/session", `{"word_length": 5, "max_guesses": 6}`).Body.Bytes(), &created)
	return "/api/session/" + created.SessionID + "/guess"
}

func TestGuessRepeatedLetters(test *testing.T) {
	router := apitest.NewRouter()

	recorder := apitest.ServeBody(router, http.MethodPost, newGuessPath(router), `{"guess": "eerie"}`)
	if recorder.Code != http.StatusOK {
		test.Fatalf("actual = %d, desired = %d", recorder.Code, http.StatusOK)
	}

	var response session.Response
	json.Unmarshal(recorder.Body.Bytes(), &response)

	desiredResult := []gamesession.LetterResult{gamesession.Absent, gamesession.Absent, gamesession.Present, gamesession.Absent, gamesession.Correct}
	if !slices.Equal(response.Guesses[0].Result, desiredResult) {
		test.Errorf("actual = %v, desired = %v", response.Guesses[0].Result, desiredResult)
	}
	if response.HiddenWord != "" {
		test.Errorf("hidden word revealed before the game is complete: %s", response.HiddenWord)
	}
}

func TestGuessWinThenGuessAgain(test *testing.T) {
	router := apitest.NewRouter()
	guessPath := newGuessPath(router)

	winRecorder := apitest.ServeBody(router, http.MethodPost, guessPath, `{"guess": "CRATE"}`)
	if winRecorder.Code != http.StatusOK {
		test.Fatalf("win: actual = %d, desired = %d", winRecorder.Code, http.StatusOK)
	}

	var response session.Response
	json.Unmarshal(winRecorder.Body.Bytes(), &response)
	if response.HiddenWord != "crate" {
		test.Errorf("win: actual hidden word = %q, desired = %q", response.HiddenWord, "crate")
	}

	againRecorder := apitest.ServeBody(router, http.MethodPost, guessPath, `{"guess": "crate"}`)
	if againRecorder.Code != http.StatusConflict {
		test.Errorf("again: actual = %d, desired = %d", againRecorder.Code, http.StatusConflict)
	}
}
