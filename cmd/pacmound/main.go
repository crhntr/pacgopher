package main

import (
	"log"
	"net/http"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
)

func main() {
	p := agents.Naive{}
	mux := pacmound.NewGameMux(&p)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
