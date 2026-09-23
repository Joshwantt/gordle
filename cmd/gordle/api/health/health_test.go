package health_test

import (
	"net/http"
	"testing"

	"github.com/Joshwantt/gordle/cmd/gordle/api/apitest"
)

func TestGet(test *testing.T) {
	recorder := apitest.ServeTest(http.MethodGet, "/api/health")

	if recorder.Code != http.StatusOK {
		test.Errorf("actual = %d, desired = %d", recorder.Code, http.StatusOK)
	}
}
