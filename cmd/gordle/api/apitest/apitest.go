package apitest

// apitest provides helpers for testing API routes.

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

func NewRouter() http.Handler {
	sessionStore := store.New()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.Register(router.Group("/api"), sessionStore, map[int][]string{5: {"crate"}}, map[int]map[string]bool{5: {"crate": true, "eerie": true}})
	return router
}

func Serve(router http.Handler, method, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}

// serve but with a body
func ServeBody(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(method, path, strings.NewReader(body)))
	return recorder
}

func ServeTest(method, path string) *httptest.ResponseRecorder {
	return Serve(NewRouter(), method, path)
}
