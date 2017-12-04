package pacmound

import "fmt"

func (m *Maze) loop() bool {
	for x := range *m {
		for y := range (*m)[x] {
			if agent := m.Occupant(x, y); agent != nil {
				agentIntent := agent.a.CalculateIntent()
				agent.previousScore = agent.score
				xIntent, yIntent := move(m, agentIntent, agent.x, agent.y)

				if agent.IsGopher() {
					fmt.Printf("(x: %d, y: %d) %s (x: %d, y: %d)\n\n\n", agent.x, agent.y, agentIntent, xIntent, yIntent)
				}

				defer func(x, y int) {
					if !(*m)[xIntent][yIntent].IsObstructed() && !(*m)[xIntent][yIntent].IsOccupied() {
						agent.x, agent.y = xIntent, yIntent
						(*m)[xIntent][yIntent].agent = (*m)[x][y].agent
						(*m)[x][y].agent = nil
					}
				}(x, y)

				if otherAgent := m.Occupant(xIntent, yIntent); otherAgent != nil {
					if fight(agent, otherAgent) {
						return false
					}
				}

				if agent.t > 0 {
					agent.score += m.getReward(x, y)
					if m.IsObsticle(xIntent, yIntent) {
						agent.Damage(DamageColision)
					}
					agent.Damage(DamageAgeing)
				}
			}
		}
	}
	return true
}
