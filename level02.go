package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Level02(gopher, python Agent) {
	loopCount, maxLoops := 0.0, 8.0*8.0

	level02(gopher, python, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || agentData.score >= (63-(loopCount*LivingCost))-0.001 || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}

func level02(gopher, python Agent, loop func(m *Maze, agentData *AgentData) bool) {
	const size = 12
	maze := NewEmptyMaze(size, size)

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if !maze[x][y].obsticle {
				maze[x][y].reward = 1
			}
		}
	}

	for i := 0; i < size; i++ {
		maze.setObsticle(0, i)
		maze.setObsticle(i, 0)
		maze.setObsticle(i, size-1)
		maze.setObsticle(size-1, i)
	}

	maze.setObsticle(5, 5)
	maze.setObsticle(5, 6)
	maze.setObsticle(5, 7)
	maze.setObsticle(5, 8)
	maze.setObsticle(5, 9)

	maze.setObsticle(5, 5)
	maze.setObsticle(5, 6)
	maze.setObsticle(6, 5)
	maze.setObsticle(6, 6)
	maze.setObsticle(7, 5)
	maze.setObsticle(7, 6)
	maze.setObsticle(8, 5)
	maze.setObsticle(8, 6)
	maze.setObsticle(9, 5)
	maze.setObsticle(9, 6)

	maze.setObsticle(9, 9)
	maze.setObsticle(9, 8)
	maze.setObsticle(8, 8)
	maze.setObsticle(8, 9)

	maze.setObsticle(2, 9)
	maze.setObsticle(3, 9)
	maze.setObsticle(2, 8)
	maze.setObsticle(3, 8)
	maze.setObsticle(2, 7)
	maze.setObsticle(3, 7)
	maze.setObsticle(2, 6)
	maze.setObsticle(3, 4)

	maze.setObsticle(4, 6)
	maze.setObsticle(4, 5)
	maze.setObsticle(4, 4)
	maze.setObsticle(4, 3)
	maze.setObsticle(4, 2)
	maze.setObsticle(4, 1)

	maze[2][2].reward = 0

	gopherData, err := maze.setAgent(2, 2, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	pythonData, err := maze.setAgent(7, 7, python)
	must(err)
	pythonData.t = -1
	python.SetScopeGetter(newScopeGetter(maze, pythonData))
	python.SetScoreGetter(pythonData.Score)

	for loop(&maze, gopherData) {
	}
}

func Level02Handler(gopher, python Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxLoops := 500
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || loopLimit > maxLoops {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := struct {
			MaxSteps int                `json:"maxSteps"`
			Scores   []float64          `json:"scores"`
			States   [][][]EncodedBlock `json:"states"`
		}{}
		data.MaxSteps = loopLimit

		level02(gopher, python, func(m *Maze, agentData *AgentData) bool {
			data.States = append(data.States, m.encodable())
			data.Scores = append(data.Scores, agentData.score)

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || loopCount > loopLimit || agentData.dead {
				return false
			}
			loopCount++
			return true
		})

		json.NewEncoder(w).Encode(data)
	}
}
