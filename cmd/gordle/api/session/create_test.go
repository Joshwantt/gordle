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

	createRecorder := apitest.Serve(router, http.MethodPost, "/api/session")
	if createRecorder.Code != http.StatusCreated {
		test.Fatalf("create: actual = %d, desired = %d", createRecorder.Code, http.StatusCreated)
	}

	var created session.Response
	json.Unmarshal(createRecorder.Body.Bytes(), &created)

	getRecorder := apitest.Serve(router, http.MethodGet, "/api/session/"+created.SessionID)
	if getRecorder.Code != http.StatusOK {
		test.Fatalf("get: actual = %d, desired = %d", getRecorder.Code, http.StatusOK)
	}
	if getRecorder.Body.String() != createRecorder.Body.String() {
		test.Errorf("get: actual = %s, desired = %s", getRecorder.Body, createRecorder.Body)
	}
}
