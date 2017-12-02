package main

import (
	"log"
	"net/http"

	"github.com/crhntr/comp469/pacman"
	"github.com/crhntr/comp469/pacman/agent"
)

func main() {
	p := agent.Naive{}
	mux := pacman.NewGameMux(&p)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
