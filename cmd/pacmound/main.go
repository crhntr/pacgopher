package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
	"github.com/crhntr/pacmound/agents/markov"
)

func getPython() pacmound.Agent {
	return &agents.Random{}
}

func main() {
	var (
		loops int
		serve bool
	)
	flag.IntVar(&loops, "loops", 0, "")
	flag.BoolVar(&serve, "serve", false, "")
	flag.Parse()

	rand.Seed(time.Now().Unix())

	agent := &markov.Agent{
		LearningRate: 0.01,
	}
	// agent := &agents.Random{}
	fmt.Println(agent)

	getGopher := func() pacmound.Agent {
		return agent
	}

	for i := 0; i < loops; i++ {
		//fmt.Printf("loop %d\n", i)
		pacmound.Level04(getGopher, getPython)
	}
	fmt.Println(agent)

	if serve {
		mux := pacmound.NewGameMux(getGopher, getPython)
		log.Println("serving on https://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}
}
