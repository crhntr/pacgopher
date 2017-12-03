package pacmound

const (
	standardReward              = 5
	standardPythonStartingScore = 1000
)

func Level00(getGopher, getPython AgentGetter) { runLevel(level00, getGopher, getPython) }
func Level01(getGopher, getPython AgentGetter) { runLevel(level01, getGopher, getPython) }
func Level02(getGopher, getPython AgentGetter) { runLevel(level02, getGopher, getPython) }
func Level03(getGopher, getPython AgentGetter) { runLevel(level03, getGopher, getPython) }
func Level04(getGopher, getPython AgentGetter) { runLevel(level04, getGopher, getPython) }

func runLevel(levelFunc LevelFunc, getGopher, getPython AgentGetter) {
	levelFunc(getGopher, getPython, func(m *Maze, agentData *AgentData) bool {
		if !m.loop() {
			return false
		}
		return true
	})
}
