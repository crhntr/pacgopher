package agents

import (
	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents/markov"
)

var randomInterfaceTest = []pacmound.Agent{
	Random{}, Simple{}, markov.Agent{}}
