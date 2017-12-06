package markov

import (
	"math"
	"math/rand"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
	"gonum.org/v1/gonum/stat"
)

const (
	numberOfActions = 4 + 1
	InitalQ         = 0.001
)

var (
	infSmall = -100000000.0
)

type Agent struct {
	LearningRate float64 `json:"learningRate"`

	CarrotWeight   float64 `json:"carrotWeight"`
	PythonWeight   float64 `json:"pythonWeight"`
	ObsticleWeight float64 `json:"obsticleWeight"`

	scope  pacmound.ScopeGetter
	reward pacmound.RewardGetter
}

func (agent *Agent) SetRewardGetter(f pacmound.RewardGetter) { agent.reward = f }
func (agent *Agent) SetScopeGetter(f pacmound.ScopeGetter)   { agent.scope = f }

func (agent *Agent) Damage(d pacmound.Damage) {
	// log.Printf("Markov Agent Took Damage %f %v", d, d.Error())
}

func (agent *Agent) CalculateIntent() pacmound.Direction {
	if agent.CarrotWeight == 0 || agent.PythonWeight == 0 {
		agent.CarrotWeight, agent.PythonWeight, agent.ObsticleWeight = 1, 1, 1
		return pacmound.Direction(rand.Intn(4) + 1)
	}
	return agent.QLearning(agent.LearningRate)
}
func (agent *Agent) Kill() {
	// fmt.Printf("KILLED (score: %f)\n", agent.reward())
	agent.QLearning(agent.LearningRate)
}
func (p *Agent) Win() {}

func (agent *Agent) QLearning(α float64) pacmound.Direction {
	r := agent.reward.LoopReward()

	actions := agents.Actions()
	d, maxScore := pacmound.DirectionNone, infSmall
	rewards := [4]float64{}
	for i, action := range actions {
		xt, yty := action.Transform()
		block := agent.scope(xt, yty)
		if block == nil {
			continue
		}
		rewards[i] = agent.CarrotWeight * agent.calulateCarrotWeight(xt, yty)
		rewards[i] += agent.PythonWeight * agent.calulatePythonWeight(xt, yty)
		rewards[i] += agent.ObsticleWeight * agent.calulateObsticleWeight(xt, yty)
		//fmt.Printf("%s (%f)\t", action, rewards[i])
		if rewards[i] > maxScore {
			maxScore = rewards[i]
			d = action
		}
	}

	if dev := stat.StdDev(rewards[:], nil); maxScore <= infSmall || dev < 0.01 {
		d = pacmound.Direction(rand.Intn(4) + 1)
	}

	cw := agent.calulateCarrotWeight(0, 0)
	pw := agent.calulatePythonWeight(0, 0)
	ow := agent.calulateObsticleWeight(0, 0)

	q := agent.CarrotWeight * cw
	q += agent.PythonWeight * pw
	q += agent.ObsticleWeight * ow

	// fmt.Printf("CW: %f, PW: %f, OW: %f, cw: %f, pw: %f, ow: %f, q: %f, r: %f\n",
	// agent.CarrotWeight, agent.PythonWeight, agent.ObsticleWeight, cw, pw, ow, q, r)

	if !math.IsNaN(cw) && !math.IsNaN(pw) && !math.IsNaN(ow) {
		agent.CarrotWeight = agent.CarrotWeight + α*(r-q)*cw
		agent.PythonWeight = agent.PythonWeight + α*(r-q)*pw
		agent.ObsticleWeight = agent.ObsticleWeight + α*(r-q)*ow
	}
	//fmt.Println(d)

	return d
}

func (agent *Agent) calulateCarrotWeight(x, y int) float64 {
	_, _, dist := agent.scope.NearestMatching(func(b *pacmound.Block) bool {
		return b != nil && b.Reward() > 0
	}, 6)
	return 1 / (dist * dist)
}
func (agent *Agent) calulatePythonWeight(x, y int) float64 {
	_, _, dist := agent.scope.NearestMatching(func(b *pacmound.Block) bool {
		return b != nil && b.IsOccupiedWithPython()
	}, 5)
	return 1 / (dist * dist)
}
func (agent *Agent) calulateObsticleWeight(x, y int) float64 {
	dist := agent.scope.AverageMatching(func(b *pacmound.Block) bool {
		return b != nil && b.IsObstructed()
	}, 4)
	return 1 / (dist * dist)
}

func maxReward(rewards ...float64) float64 {
	maxReward := rewards[0]
	for i := 1; i < len(rewards); i++ {
		if rewards[i] > maxReward && !math.IsNaN(rewards[i]) {
			maxReward = rewards[i]
		}
	}
	return maxReward
}

func maxDirection(rewards []float64, actions []pacmound.Direction) pacmound.Direction {
	maxDir, maxReward := 0, rewards[0]
	for i := 1; i < len(rewards); i++ {
		if rewards[i] > maxReward && !math.IsNaN(rewards[i]) {
			maxDir, maxReward = i, rewards[i]
		}
	}
	return actions[maxDir]
}

func distance(xFinal, yFinal, xInital, yInital int) float64 {
	return math.Sqrt(float64((xInital-xFinal)*(xInital-xFinal) + (yInital-yFinal)*(yInital-yFinal)))
}
