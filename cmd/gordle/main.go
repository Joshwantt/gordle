package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

//go:embed words
var wordFiles embed.FS

// the whitespace separated words in the embedded file fileName
func readWords(fileName string) []string {
	contents, readError := wordFiles.ReadFile(fileName)
	if readError != nil {
		log.Fatal(readError)
	}
	return strings.Fields(string(contents))
}

// every word from wordLists as a set so checking a guess is a fast map lookup
func wordSet(wordLists ...[]string) map[string]bool {
	words := map[string]bool{}
	for _, wordList := range wordLists {
		for _, word := range wordList {
			words[word] = true
		}
	}
	return words
}

func main() {
	sessionStore := store.New()
	go sessionStore.DeleteInactiveTicker(time.Minute, time.Hour*24)

	address := flag.String("address", ":8080", "listen address")
	staticDirectory := flag.String("static", "web/dist", "directory of the built site to serve")
	flag.Parse()

	answersByLength := map[int][]string{
		3:  readWords("words/valid_answers_3"),
		4:  readWords("words/valid_answers_4"),
		5:  readWords("words/valid_answers_5"),
		6:  readWords("words/valid_answers_6"),
		7:  readWords("words/valid_answers_7"),
		8:  readWords("words/valid_answers_8"),
		9:  readWords("words/valid_answers_9"),
		10: readWords("words/valid_answers_10"),
		11: readWords("words/valid_answers_11"),
		12: readWords("words/valid_answers_12"),
		13: readWords("words/valid_answers_13"),
		14: readWords("words/valid_answers_14"),
		15: readWords("words/valid_answers_15"),
		16: readWords("words/valid_answers_16"),
	}

	// every answer is also a valid guess so join them
	guessesByLength := map[int]map[string]bool{
		3:  wordSet(readWords("words/valid_guess_3"), answersByLength[3]),
		4:  wordSet(readWords("words/valid_guess_4"), answersByLength[4]),
		5:  wordSet(readWords("words/valid_guess_5"), answersByLength[5]),
		6:  wordSet(readWords("words/valid_guess_6"), answersByLength[6]),
		7:  wordSet(readWords("words/valid_guess_7"), answersByLength[7]),
		8:  wordSet(readWords("words/valid_guess_8"), answersByLength[8]),
		9:  wordSet(readWords("words/valid_guess_9"), answersByLength[9]),
		10: wordSet(readWords("words/valid_guess_10"), answersByLength[10]),
		11: wordSet(readWords("words/valid_guess_11"), answersByLength[11]),
		12: wordSet(readWords("words/valid_guess_12"), answersByLength[12]),
		13: wordSet(readWords("words/valid_guess_13"), answersByLength[13]),
		14: wordSet(readWords("words/valid_guess_14"), answersByLength[14]),
		15: wordSet(readWords("words/valid_guess_15"), answersByLength[15]),
		16: wordSet(readWords("words/valid_guess_16"), answersByLength[16]),
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
