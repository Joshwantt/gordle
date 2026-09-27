package session_id_test

import (
	"net/http"
	"testing"
	"uuid"

	"github.com/Joshwantt/gordle/cmd/gordle/api/apitest"
)

func TestGetInvalidID(test *testing.T) {
	recorder := apitest.ServeTest(http.MethodGet, "/api/session/not-a-uuid")

	if recorder.Code != http.StatusBadRequest {
		test.Errorf("actual = %d, desired = %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetUnknownSession(test *testing.T) {
	recorder := apitest.ServeTest(http.MethodGet, "/api/session/"+uuid.New().String())

	if recorder.Code != http.StatusNotFound {
		test.Errorf("actual = %d, desired = %d", recorder.Code, http.StatusNotFound)
	}
}
