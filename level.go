package pacmound

const (
	standardReward = 5
)

func Level00(getGopher, getPython AgentGetter) { runLevel(level00, getGopher, getPython, 1000) }
func Level01(getGopher, getPython AgentGetter) { runLevel(level01, getGopher, getPython, 1000) }
func Level02(getGopher, getPython AgentGetter) { runLevel(level02, getGopher, getPython, 1000) }
func Level03(getGopher, getPython AgentGetter) { runLevel(level03, getGopher, getPython, 1000) }
func Level04(getGopher, getPython AgentGetter) { runLevel(level04, getGopher, getPython, 1000) }

func runLevel(levelFunc LevelFunc, getGopher, getPython AgentGetter, maxLoops int) {
	loopCount := 0
	levelFunc(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() || loopCount > maxLoops {
			return false
		}
		loopCount++
		return true
	})
}
