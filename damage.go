package pacmound

import "fmt"

type Damage float64

const (
	DamageAgeing    Damage = 0.5
	DamageLostFight Damage = 500.0
	DamageColision  Damage = 30

	DeathPenatlty = 20
)

func (d Damage) Error() string {
	switch d {
	case DamageAgeing:
		return "Damage due to ageing"
	case DamageLostFight:
		return "Damage due to fight lost"
	case DamageColision:
		return "Damage due to colision"
	default:
		return fmt.Sprintf("unknown damage %f", d)
	}
}

func (ad *AgentData) Damage(d Damage) {
	ad.score -= float64(d)
	ad.a.Damage(d)
	if ad.dead = ad.score <= 0; ad.dead {
		ad.score -= DeathPenatlty
		ad.a.Kill()
	}
}
