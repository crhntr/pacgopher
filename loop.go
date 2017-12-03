package pacmound

func (m *Maze) loop() bool {
	gopherDied := false
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
					fight(agent, otherAgent)
					if agent.IsGopher() && agent.dead || otherAgent.IsGopher() && otherAgent.dead {
						return false
					}
				}

				if agent.t > 0 {
					defer func(x, y int) {
						agent.score += (*m)[x][y].reward
						(*m)[x][y].reward = 0
					}(x, y)
					agent.score -= LivingCost
				}

				defer func(x, y int, hitObsticle bool) {
					if !hitObsticle && !gopherDied {
						agent.x = xIntent
						agent.y = yIntent
						(*m)[xIntent][yIntent].agent = (*m)[x][y].agent
						(*m)[x][y].agent = nil
					}
				}(x, y, hitObsticle)
			}
		}
	}
	return true
}
