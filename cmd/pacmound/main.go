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
		CarrotWeight:   1,
		ObsticleWeight: 1,
		PythonWeight:   1,

		DiscountFactor: 0.1,
		LearningRate:   0.1,
	}
	getGopher := func() pacmound.Agent {
		return agent
	}

	mux := pacmound.NewGameMux(getGopher, getPython)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
