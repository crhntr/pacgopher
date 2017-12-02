package pacmound

import (
	"math"
)

func (m *Maze) loop() bool {
	for x := range *m {
		for y := range (*m)[x] {

			if agent := m.Occupant(x, y); agent != nil {
				agentIntent := agent.a.CalculateIntent()
				xIntent, yIntent := move(m, agentIntent, agent.x, agent.y)
				hitObsticle := m.IsObsticle(xIntent, yIntent)
				if hitObsticle {
					agent.score -= ObsticleCollisionCost
				}

				if otherAgent := m.Occupant(xIntent, yIntent); otherAgent != nil {
					if math.Abs(otherAgent.score) > math.Abs(agent.score) {
						agent.RecordKill()
					} else {
						otherAgent.RecordKill()
					}
				}

				points, err := m.RewardAt(x, y)
				if err != nil {
					agent.a.Warning(err)
				}
				agent.score += points

				defer func(x, y int) {
					if !hitObsticle {
						agent.x = xIntent
						agent.y = yIntent
						(*m)[xIntent][yIntent].agent = (*m)[x][y].agent
						(*m)[x][y].agent = nil
					}
				}(x, y)
			}
		}

		for x := range *m {
			for y := range (*m)[x] {
				if agent := (*m)[x][y].agent; agent != nil {
					if agent.t > 0 {
						(*m)[x][y].reward = 0
						agent.score -= LivingCost
						if agent.dead {
							return false
						}
					}
				}
			}
		}
	}
	return true
}
