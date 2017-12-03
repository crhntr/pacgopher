package pacmound

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

func level04(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool) {
	const height, width = 8, 32
	maze := NewEmptyMaze(height, width)
	for x := 0; x < height; x++ {
		maze.setObsticle(x, 0)
		maze.setObsticle(x, width-1)
		for y := 0; y < width; y++ {
			maze.setObsticle(0, y)
			maze.setObsticle(height-1, y)
		}
	}

	for x := 0; x < height-1; x++ {
		for y := 0; y < width-1; y++ {
			if !maze[x][y].obsticle {
				if rand.Intn(100) > 97 {
					python := getPython()
					pythonData, err := maze.setAgent(x, y, python)
					must(err)
					pythonData.t = -1
					pythonData.score = DeathCost
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

	maze[2][2].reward = 0
	maze[2][2].obsticle = false
	gopher := getGopher()
	gopherData, err := maze.setAgent(2, 2, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	for loop(&maze, gopherData) {
	}
}

func Level04Handler(getGopher, getPython AgentGetter, mut *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mut.Lock()
		defer mut.Unlock()

		maxLoops := MaxLoops
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = loopLimit

		level04(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
			data.States = append(data.States, m.encodable())
			data.Scores = append(data.Scores, agentData.score)

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || loopCount > loopLimit || agentData.dead {
				data.Scores = append(data.Scores, agentData.score)
				return false
			}
			loopCount++
			return true
		})

		data.Agent = getGopher()
		json.NewEncoder(w).Encode(data)
	}
}
