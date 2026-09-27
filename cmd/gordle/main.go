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

func main() {
	sessionStore := store.New()
	go sessionStore.DeleteInactiveTicker(time.Minute, time.Hour*24)

	address := flag.String("address", ":8080", "listen address")
	staticDirectory := flag.String("static", "web/dist", "directory of the built site to serve")
	flag.Parse()

	answersByLength := map[int][]string{
		5: strings.Fields(validAnswers5),
	}

	router := gin.Default()

	api.Register(router.Group("/api"), sessionStore, answersByLength)

	router.NoRoute(gin.WrapH(http.FileServer(http.Dir(*staticDirectory))))

	log.Fatal(router.Run(*address))
}
