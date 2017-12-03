package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Level01(gopher, python Agent) {
	loopCount, maxLoops := 0.0, 2000.0

	level01(gopher, python, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || agentData.score >= (63-(loopCount*LivingCost))-0.001 || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}

func level01(gopher, python Agent, loop func(m *Maze, agentData *AgentData) bool) {
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

	gopherData, err := maze.setAgent(2, 2, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	pythonData, err := maze.setAgent(7, 7, python)
	must(err)
	pythonData.t = -1
	pythonData.score = DeathCost
	python.SetScopeGetter(newScopeGetter(maze, pythonData))
	python.SetScoreGetter(pythonData.Score)

	for loop(&maze, gopherData) {
	}
}

func Level01Handler(getGopher, getPython AgentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		training := r.URL.Query().Get("train") == trueStr

		gopher, python := getGopher(), getPython()
		maxLoops := MaxLoops
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || loopLimit > maxLoops {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = loopLimit

		level01(gopher, python, func(m *Maze, agentData *AgentData) bool {
			if !training {
				data.States = append(data.States, m.encodable())
				data.Scores = append(data.Scores, agentData.score)
			}

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || (!training && loopCount > loopLimit) || agentData.dead {
				data.Scores = append(data.Scores, agentData.score)
				gopher.CalculateIntent()
				return false
			}
			loopCount++
			return true
		})

		data.Agent = gopher
		json.NewEncoder(w).Encode(data)
	}
}
