package pacmound

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

func Level03(gopher, python Agent) {
	loopCount, maxLoops := 0.0, 8.0*8.0

	level03(gopher, python, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || agentData.score >= (63-(loopCount*LivingCost))-0.001 || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}

func level03(gopher, python Agent, loop func(m *Maze, agentData *AgentData) bool) {
	const height, width = 45, 50
	maze := NewEmptyMaze(height, width)
	for x := 0; x < height; x++ {
		maze.setObsticle(x, 0)
		maze.setObsticle(x, width-1)
		for y := 0; y < width; y++ {
			maze.setObsticle(0, y)
			maze.setObsticle(height-1, y)
		}
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			if !maze[x][y].obsticle {
				if rand.Intn(100) > 93 {
					pythonData, err := maze.setAgent(7, 7, python)
					must(err)
					pythonData.t = -1
					pythonData.score = 1000
					python.SetScopeGetter(newScopeGetter(maze, pythonData))
					python.SetScoreGetter(pythonData.Score)
				} else if rand.Intn(100) < 5 {
					maze.setObsticle(x, y)
				} else if rand.Intn(100) > 100-30 {
					maze.setReward(x, y, float64(int64(rand.Float64()*10*100))/100)
				}
			}
		}
	}

	maze[1][1].reward = 0
	gopherData, err := maze.setAgent(1, 1, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	for loop(&maze, gopherData) {
	}
}

func Level03Handler(gopher, python Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxLoops := 5
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

		level03(gopher, python, func(m *Maze, agentData *AgentData) bool {
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
