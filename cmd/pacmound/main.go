package main

import (
	"log"
	"net/http"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
)

func main() {
	gopher := agents.Naive{}
	python1, python2, python3 := agents.Ghost{}, agents.Ghost{}, agents.Ghost{}
	mux := pacmound.NewGameMux(&gopher, &python1, &python2, &python3)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
