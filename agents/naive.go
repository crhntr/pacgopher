package agents

import (
	"math/rand"

	"github.com/crhntr/pacmound"
)

type Naive struct {
	// dead  bool
	score pacmound.ScoreGetter
	scope pacmound.ScopeGetter

	warning error
}

func (p *Naive) SetScoreGetter(f pacmound.ScoreGetter) { p.score = f }
func (p *Naive) SetScopeGetter(f pacmound.ScopeGetter) { p.scope = f }
func (p *Naive) Warning(err error)                     { p.warning = err }

func (p *Naive) CalculateIntent() pacmound.Direction {
	// time.Sleep(time.Second / 10)

	directions := directionsSlice()
	rewards := make([]float64, len(directions))

	for i, dir := range directions {
		x, y := dir.Transform()

		dirReward := 0.0
		out := 1
		for {
			b := p.scope(x*out, y*out)
			if b.IsOccupied() {
				dirReward *= -1000
			}
			if b == nil || b.IsObstructed() || out > 5 {
				break
			}
			dirReward += b.Reward()
			out++
		}

		rewards[i] = dirReward
	}
	_, directions = removeMinimumScoringDirections(rewards, directions)
	return directions[rand.Intn(len(directions))]
}

// func (p *Naive) Kill()                                 { p.dead = true }
