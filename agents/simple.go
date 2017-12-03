package agents

import (
	"log"
	"math/rand"

	"github.com/crhntr/pacmound"
)

type Simple struct {
	dead  bool
	score pacmound.ScoreGetter
	scope pacmound.ScopeGetter

	warning error
}

func (p *Simple) SetScoreGetter(f pacmound.ScoreGetter) { p.score = f }
func (p *Simple) SetScopeGetter(f pacmound.ScopeGetter) { p.scope = f }
func (p *Simple) Warning(err error)                     { p.warning = err }

func (p *Simple) CalculateIntent() pacmound.Direction {
	// time.Sleep(time.Second / 10)

	directions := Actions()
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

func (p *Simple) Kill()                    { p.dead = true }
func (p *Simple) Damage(d pacmound.Damage) { log.Printf("Simple took damage: %s", d) }
