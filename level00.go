package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Level00(agent Agent) {
	loopCount, maxLoops := 0, 1000

	level00(agent, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || agentData.score >= (63-(float64(loopCount)*LivingCost))-0.001 || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}

func level00(player Agent, loop func(m *Maze, agentData *AgentData) bool) {
	const size = 10
	maze := NewEmptyMaze(size, size)
	for i := 0; i < size; i++ {
		maze.setObsticle(0, i)
		maze.setObsticle(i, 0)
		maze.setObsticle(i, size-1)
		maze.setObsticle(size-1, i)
	}

	for x := 1; x < size-1; x++ {
		for y := 1; y < size-1; y++ {
			maze[x][y].reward = 1
		}
	}
	maze[2][2].reward = 0

	agentData, err := maze.setAgent(2, 2, player)
	must(err)
	agentData.t = 1
	player.SetScopeGetter(newScopeGetter(maze, agentData))
	player.SetScoreGetter(agentData.Score)

	for loop(&maze, agentData) {
	}
}

func Level00Handler(agent Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		training := r.URL.Query().Get("train") == trueStr

		maxLoops := MaxLoops
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || loopLimit > maxLoops {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = loopLimit

		level00(agent, func(m *Maze, agentData *AgentData) bool {
			if !training {
				data.States = append(data.States, m.encodable())
				data.Scores = append(data.Scores, agentData.score)
			}
			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || (!training && loopCount > loopLimit) || agentData.dead {
				data.Scores = append(data.Scores, agentData.score)
				agent.CalculateIntent()
				return false
			}
			loopCount++
			return true
		})

		data.Agent = agent
		json.NewEncoder(w).Encode(data)
	}
}
