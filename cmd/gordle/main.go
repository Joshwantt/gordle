package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
	"github.com/Joshwantt/gordle/cmd/gordle/store"
)

func main() {
	sessionStore := store.New()
	go sessionStore.DeleteInactiveTicker(time.Minute, time.Hour*24)

	address := flag.String("address", ":8080", "listen address")
	staticDirectory := flag.String("static", "web/dist", "directory of the built site to serve")
	flag.Parse()

	router := gin.Default()

	api.Register(router.Group("/api"), sessionStore)

	router.NoRoute(gin.WrapH(http.FileServer(http.Dir(*staticDirectory))))

	log.Fatal(router.Run(*address))
}
