package pacmound

func level03(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool) {
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

	maze.setObsticle(5, 7)
	maze.setObsticle(5, 8)
	maze.setObsticle(5, 9)

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
	// maze.setObsticle(4, 5)
	maze.setObsticle(4, 4)
	maze.setObsticle(4, 3)
	maze.setObsticle(4, 2)
	maze.setObsticle(4, 1)

	maze.setObsticle(9, 2)
	maze.setObsticle(9, 3)
	maze.setObsticle(9, 4)

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
