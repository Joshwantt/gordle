package apitest

// apitest provides helpers for testing API routes.

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

func ServeTest(method, path string) *httptest.ResponseRecorder {
	sessionStore := store.New()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.Register(router.Group("/api"), sessionStore)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}
