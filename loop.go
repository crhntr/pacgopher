package pacman

import (
	"math"
)

func (m *Maze) loop() bool {
	for x := range *m {
		for y := range (*m)[x] {

			if agent := m.Occupant(x, y); agent != nil {
				agentIntent := agent.a.CalculateIntent()
				xIntent, yIntent := move(m, agentIntent, agent.x, agent.y)
				if m.IsObsticle(xIntent, yIntent) {
					agent.score -= ObsticleCollisionCost
					continue
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

				defer func(x, y int) {
					agent.x = xIntent
					agent.y = yIntent
					(*m)[xIntent][yIntent].agent = (*m)[x][y].agent
					(*m)[x][y].agent = nil
				}(x, y)

				if agent.t > 0 {
					if agent.score < 0 {
						return false
					}
					agent.score += points
					(*m)[x][y].reward = 0

					agent.score -= LivingCost
				}
			}
		}
	}
	return true
}
