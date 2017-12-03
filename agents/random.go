package agents

import (
	"math/rand"

	"github.com/crhntr/pacmound"
)

type Random struct{}

func (p *Random) SetScoreGetter(f pacmound.ScoreGetter) {}
func (p *Random) SetScopeGetter(f pacmound.ScopeGetter) {}
func (p *Random) Damage(d pacmound.Damage)              {}
func (p *Random) CalculateIntent() pacmound.Direction {
	return pacmound.Direction(rand.Intn(4) + 1)
}
func (p *Random) Kill() { /* What is Dead May Never Die*/ }
