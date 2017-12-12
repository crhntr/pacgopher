# pacmound
A pacman inspired package for experimenting with various AI agents.

## The "Improved" Markov Agent
The Agent found in the Markov package (pacmound/agents/markov) is based on
a template agent written by Dr. Claveau for a Intro to Artificial
Intelligence class project. I wrote similar game logic in and translated
some of the agent logic from Python to Go.

I made a few changes to the algorithm that seem to improve it,
(see comments in the code)
```go
func (agent *Agent) QLearning(α float64) pacmound.Direction {
	r := agent.reward.LoopReward()

	actions := agents.Actions()
	d, maxScore := pacmound.DirectionNone, infSmall
	rewards := [4]float64{}

	for i, action := range actions {
		xt, yt := action.Transform()
		block := agent.scope(xt, yt)
		if block == nil {
			continue
		}
		rewards[i] = agent.CarrotWeight * agent.calulateCarrotWeight(xt, yt)
		rewards[i] += agent.PythonWeight * agent.calulatePythonWeight(xt, yt

		// I also added this calulateObsticleWeight so that the agent must
		// also learn about crashing into walls and the costs associated.
		rewards[i] += agent.ObsticleWeight * agent.calulateObsticleWeight(xt, yt)
		if rewards[i] > maxScore {
			maxScore = rewards[i]
			d = action
		}
	}

	// Here I ensured that if the agent does not have a strong opinion
	// on where to go, it chooses randomly. This ensures that if the
	// agent isn't sure, it won't just go somewhere based on insufficient
	// experience and that it will go someplace based on some learning.
	if dev := stat.StdDev(rewards[:], nil); maxScore <= infSmall || dev < 0.01 {
		d = pacmound.Direction(rand.Intn(4) + 1)
	}

	cw := agent.calulateCarrotWeight(0, 0)
	pw := agent.calulatePythonWeight(0, 0)
	ow := agent.calulateObsticleWeight(0, 0)

	q := agent.CarrotWeight * cw
	q += agent.PythonWeight * pw
	q += agent.ObsticleWeight * ow

	if !math.IsNaN(cw) && !math.IsNaN(pw) && !math.IsNaN(ow) {
		agent.CarrotWeight = agent.CarrotWeight + α*(r-q)*cw
		agent.PythonWeight = agent.PythonWeight + α*(r-q)*pw
		agent.ObsticleWeight = agent.ObsticleWeight + α*(r-q)*ow
	}

	return d
}
```
This code loosely was based on:
- http://faculty.csuci.edu/David.Claveau/COMP469F17/pacman.txt
- http://faculty.csuci.edu/David.Claveau/COMP469F17/GridMDP_PM.txt
- http://faculty.csuci.edu/David.Claveau/COMP469F17/MDP4.txt
- http://faculty.csuci.edu/David.Claveau/COMP469F17/GridMDP4.txt
