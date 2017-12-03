package agents

import (
	"log"
	"math/rand"

	"github.com/crhntr/pacmound"
)

type Random struct {
	dead  bool
	score pacmound.ScoreGetter
	scope pacmound.ScopeGetter

	warning error
}

func (p *Random) SetScoreGetter(f pacmound.ScoreGetter) { p.score = f }
func (p *Random) SetScopeGetter(f pacmound.ScopeGetter) { p.scope = f }
func (p *Random) Warning(err error)                     { log.Printf("GHOST: %s", err) }
func (p *Random) CalculateIntent() pacmound.Direction {
	return directionsSlice()[rand.Intn(4)]
}

// func (p *Random) Kill() { /* What is Dead May Never Die*/ }
