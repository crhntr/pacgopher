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
			maze[x][y].reward = 1
		}
	}
	maze[2][2].reward = 0

	gopher := getGopher()
	gopherData, err := maze.setAgent(2, 2, gopher)
	must(err)
	gopherData.t = 1
	gopher.SetScopeGetter(newScopeGetter(maze, gopherData))
	gopher.SetScoreGetter(gopherData.Score)

	python := getPython()
	pythonData, err := maze.setAgent(7, 7, python)
	must(err)
	pythonData.t = -1
	pythonData.score = DeathCost
	python.SetScopeGetter(newScopeGetter(maze, pythonData))
	python.SetScoreGetter(pythonData.Score)

	for loop(&maze, gopherData) {
	}
}
