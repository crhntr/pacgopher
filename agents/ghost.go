package agents

import (
	"log"
	"math/rand"

	"github.com/crhntr/pacmound"
)

type Ghost struct {
	dead  bool
	score pacmound.ScoreGetter
	scope pacmound.ScopeGetter

	warning error
}

func (p *Ghost) SetScoreGetter(f pacmound.ScoreGetter) { p.score = f }
func (p *Ghost) SetScopeGetter(f pacmound.ScopeGetter) { p.scope = f }
func (p *Ghost) Warning(err error)                     { log.Printf("GHOST: %s", err) }
func (p *Ghost) CalculateIntent() pacmound.Direction {
	return directions[rand.Intn(len(directions))]
}

// func (p *Ghost) Kill() { /* What is Dead May Never Die*/ }
