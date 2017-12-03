package pacmound

func (m *Maze) loop() bool {
	for x := range *m {
		for y := range (*m)[x] {
			if agent := m.Occupant(x, y); agent != nil {
				agentIntent := agent.a.CalculateIntent()
				xIntent, yIntent := move(m, agentIntent, agent.x, agent.y)

				if otherAgent := m.Occupant(xIntent, yIntent); otherAgent != nil {
					if fight(agent, otherAgent) {
						return false
					}
				}

				hitObsticle := m.IsObsticle(xIntent, yIntent)

				if agent.t > 0 {
					agent.score += (*m)[x][y].reward
					(*m)[x][y].reward = 0

					if hitObsticle {
						agent.score -= ObsticleCollisionCost
					}

					if agent.score < 0 {
						agent.RecordKill()
						return false
					}

					agent.score -= LivingCost
				}

				defer func(x, y int) {
					if !(*m)[xIntent][yIntent].IsObstructed() && !(*m)[xIntent][yIntent].IsOccupied() {
						agent.x, agent.y = xIntent, yIntent
						(*m)[xIntent][yIntent].agent = (*m)[x][y].agent
						(*m)[x][y].agent = nil
					}
				}(x, y)
			}
		}
	}
	return true
}
