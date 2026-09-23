package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Joshwantt/gordle/cmd/gordle/api"
)

func main() {
	address := flag.String("address", ":8080", "listen address")
	staticDirectory := flag.String("static", "web/dist", "directory of the built site to serve")
	flag.Parse()

	router := gin.Default()

	api.Register(router.Group("/api"))

	router.NoRoute(gin.WrapH(http.FileServer(http.Dir(*staticDirectory))))

	log.Fatal(router.Run(*address))
}
