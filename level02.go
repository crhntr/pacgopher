package pacmound

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Level02(gopher, python1, python2 Agent) {
	loopCount, maxLoops := 0.0, 8.0*8.0

	level02(gopher, python1, python2, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || agentData.score >= (63-(loopCount*LivingCost))-0.001 || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}

func level02(gopher, python1, python2 Agent, loop func(m *Maze, agentData *AgentData) bool) {
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
				maze.setReward(x, y, 1)
			}
		}
	}

	maze[1][1].reward = 0
	maze[1][1].obsticle = false
	gopherData, err := maze.setAgent(1, 1, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	python1Data, err := maze.setAgent(3, 8, python1)
	must(err)
	python1Data.t = -1
	python1Data.score = 1000
	python1.SetScopeGetter(newScopeGetter(maze, python1Data))
	python1.SetScoreGetter(python1Data.Score)

	python2Data, err := maze.setAgent(6, 5, python2)
	must(err)
	python2Data.t = -1
	python2Data.score = 1000
	python2.SetScopeGetter(newScopeGetter(maze, python2Data))
	python2.SetScoreGetter(python2Data.Score)

	for loop(&maze, gopherData) {
	}
}

func Level02Handler(getGopher, getPython AgentGetter) http.HandlerFunc {
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

		gopher := getGopher()
		level02(gopher, getPython(), getPython(), func(m *Maze, agentData *AgentData) bool {
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
