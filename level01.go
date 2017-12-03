package pacmound

func level01(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool) {
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
			maze[x][y].reward = standardReward
		}
	}

	gopher := getGopher()
	gopherData, err := maze.setAgent(2, 2, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	python1 := getPython()
	python1Data, err := maze.setAgent(5, 7, python1)
	must(err)
	python1Data.t = -1
	python1Data.score = standardPythonStartingScore
	python1.SetScopeGetter(newScopeGetter(maze, python1Data))
	python1.SetScoreGetter(python1Data.Score)

	python2 := getPython()
	python2Data, err := maze.setAgent(7, 5, python2)
	must(err)
	python2Data.t = -1
	python2Data.score = standardPythonStartingScore
	python2.SetScopeGetter(newScopeGetter(maze, python2Data))
	python2.SetScoreGetter(python2Data.Score)

	for loop(&maze, gopherData) {
	}
}
