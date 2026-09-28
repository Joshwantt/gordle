package session_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Joshwantt/gordle/cmd/gordle/api/apitest"
	"github.com/Joshwantt/gordle/cmd/gordle/api/session"
)

func TestCreateThenGet(test *testing.T) {
	router := apitest.NewRouter()

	createRecorder := apitest.ServeBody(router, http.MethodPost, "/api/session", `{"word_length": 5, "max_guesses": 3}`)
	if createRecorder.Code != http.StatusCreated {
		test.Fatalf("create: actual = %d, desired = %d", createRecorder.Code, http.StatusCreated)
	}

	var created session.Response
	json.Unmarshal(createRecorder.Body.Bytes(), &created)
	if created.WordLength != 5 || created.MaxGuesses != 3 {
		test.Errorf("create: actual word_length = %d, max_guesses = %d, desired = 5, 3", created.WordLength, created.MaxGuesses)
	}

	getRecorder := apitest.Serve(router, http.MethodGet, "/api/session/"+created.SessionID)
	if getRecorder.Code != http.StatusOK {
		test.Fatalf("get: actual = %d, desired = %d", getRecorder.Code, http.StatusOK)
	}
	if getRecorder.Body.String() != createRecorder.Body.String() {
		test.Errorf("get: actual = %s, desired = %s", getRecorder.Body, createRecorder.Body)
	}
}

func TestCreateInvalidRequest(test *testing.T) {
	router := apitest.NewRouter()

	for _, body := range []string{
		`{}`,
		`{"word_length": 5}`,
		`{"max_guesses": 6}`,
		`{"word_length": 2, "max_guesses": 6}`,
		`{"word_length": 17, "max_guesses": 6}`,
		`{"word_length": 5, "max_guesses": 0}`,
		`{"word_length": 5, "max_guesses": 21}`,
	} {
		recorder := apitest.ServeBody(router, http.MethodPost, "/api/session", body)
		if recorder.Code != http.StatusBadRequest {
			test.Errorf("%s: actual = %d, desired = %d", body, recorder.Code, http.StatusBadRequest)
		}
	}
}
