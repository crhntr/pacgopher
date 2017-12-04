package pacmound

func level00(getGopher, getPython AgentGetter, loop func(m *Maze, agentData *AgentData) bool) {
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
			maze.setReward(x, y, standardReward)
		}
	}

	player := getGopher()
	agentData, err := maze.setAgent(2, 2, player)
	must(err)
	agentData.t = 1
	player.SetScopeGetter(newScopeGetter(maze, agentData))
	player.SetRewardGetter(agentData)

	for loop(&maze, agentData) {
	}
}
