package main

import (
	"flag"
	"fmt"
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
	var loops int
	flag.IntVar(&loops, "loops", 1, "")
	flag.Parse()

	agent := &markov.Agent{
		LearningRate: 0.07,
	}
	fmt.Println(agent)

	getGopher := func() pacmound.Agent {
		return agent
	}

	for i := 0; i < loops; i++ {
		log.Printf("\n\n\n\n\nLOOP %d\n\n\n\n\n", i)
		pacmound.Level01(getGopher, getPython)
	}
	fmt.Println(agent)

	mux := pacmound.NewGameMux(getGopher, getPython)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
