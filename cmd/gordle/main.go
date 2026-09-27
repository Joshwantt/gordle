package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

//go:embed valid_answers_5
var validAnswers5 string

//go:embed valid_guess_5
var validGuess5 string

func main() {
	sessionStore := store.New()
	go sessionStore.DeleteInactiveTicker(time.Minute, time.Hour*24)

	address := flag.String("address", ":8080", "listen address")
	staticDirectory := flag.String("static", "web/dist", "directory of the built site to serve")
	flag.Parse()

	answersByLength := map[int][]string{
		5: strings.Fields(validAnswers5),
	}

	// every answer is also a valid guess so join them.
	// store these a set so checking a guess is a fast map lookup
	guessesByLength := map[int]map[string]bool{5: {}}
	for _, guess := range append(strings.Fields(validGuess5), answersByLength[5]...) {
		guessesByLength[5][guess] = true
	}

	router := gin.Default()

	api.Register(router.Group("/api"), sessionStore, answersByLength, guessesByLength)

	router.NoRoute(gin.WrapH(http.FileServer(http.Dir(*staticDirectory))))

	// timeouts so slow or idle clients can't hold connections open forever
	server := &http.Server{
		Addr:              *address,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	log.Fatal(server.ListenAndServe())
}
