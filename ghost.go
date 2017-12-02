package pacman

import (
	"log"
	"math/rand"
)

type Ghost struct {
	dead  bool
	score ScoreGetter
	scope ScopeGetter

	warning error
}

func (p *Ghost) SetScoreGetter(f ScoreGetter) { p.score = f }
func (p *Ghost) SetScopeGetter(f ScopeGetter) { p.scope = f }
func (p *Ghost) Warning(err error)            { log.Printf("GHOST: %s", err) }
func (p *Ghost) CalculateIntent() Direction {
	return directions[rand.Intn(len(directions))]
}
func (p *Ghost) Kill() { /* What is Dead May Never Die*/ }
