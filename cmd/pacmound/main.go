package main

import (
	"log"
	"net/http"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
	"github.com/crhntr/pacmound/agents/markov"
)

func getPython() pacmound.Agent {
	return &agents.Random{}
}

func main() {
	agent := &markov.Agent{
		DiscountFactor: 0.1,
		LearningRate:   0.5,
		QTable:         markov.NewQTable(5, 5),
	}
	getGopher := func() pacmound.Agent {
		return agent
	}

	mux := pacmound.NewGameMux(getGopher, getPython)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
