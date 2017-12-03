package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

func level02(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool) {
	const height, width = 9, 11
	maze := NewEmptyMaze(height, width)
	for x := 0; x < height; x++ {
		maze.setObsticle(x, 0)
		maze.setObsticle(x, width-1)
		for y := 0; y < width; y++ {
			maze.setObsticle(0, y)
			maze.setObsticle(height-1, y)
		}
	}

	for x := 1; x < height-1; x++ {
		for y := 1; y < width-1; y++ {
			if (y+2)%2 == 0 && (x+2)%2 == 0 {
				maze.setObsticle(x, y)
			} else {
				maze.setReward(x, y, standardReward)
			}
		}
	}

	maze[1][1].obsticle = false
	gopher := getGopher()
	gopherData, err := maze.setAgent(1, 1, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	python1 := getPython()
	python1Data, err := maze.setAgent(3, 8, python1)
	must(err)
	python1Data.t = -1
	python1Data.score = standardPythonStartingScore
	python1.SetScopeGetter(newScopeGetter(maze, python1Data))
	python1.SetScoreGetter(python1Data.Score)

	python2 := getPython()
	python2Data, err := maze.setAgent(6, 5, python2)
	must(err)
	python2Data.t = -1
	python2Data.score = standardPythonStartingScore
	python2.SetScopeGetter(newScopeGetter(maze, python2Data))
	python2.SetScoreGetter(python2Data.Score)

	for loop(&maze, gopherData) {
	}
}

func Level02Handler(getGopher, getPython AgentGetter, mut *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mut.Lock()
		defer mut.Unlock()

		maxLoops := MaxLoops
		loopLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || loopLimit > maxLoops {
			loopLimit = maxLoops
		}
		loopCount := 0

		data := LevelData{}
		data.MaxSteps = loopLimit

		gopher := getGopher()
		level02(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
			data.States = append(data.States, m.encodable())
			data.Scores = append(data.Scores, agentData.score)

			remReward := m.RemainingReward()

			if !m.loop() || remReward <= 0 || loopCount > loopLimit || agentData.dead {
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
