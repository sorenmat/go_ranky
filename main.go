package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/sorenmat/ranky/matchservice"
	"github.com/sorenmat/ranky/playerservice"
)

func main() {

	restful.Add(playerservice.New())
	restful.Add(matchservice.New())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
