package markov

import (
	"log"
	"math"
	"math/rand"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
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
	reward pacmound.ScoreGetter

	prevousAction pacmound.Direction

	previousRewardTotal float64
}

func (agent *Agent) SetScoreGetter(f pacmound.ScoreGetter) { agent.reward = f }
func (agent *Agent) SetScopeGetter(f pacmound.ScopeGetter) { agent.scope = f }

func (agent *Agent) Damage(d pacmound.Damage) {
	log.Printf("Markov Agent Took Damage %f %v", d, d.Error())
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
	totalReward := agent.reward()
	r := totalReward - agent.previousRewardTotal
	agent.previousRewardTotal = totalReward

	actions := agents.Actions()
	d, maxScore := pacmound.DirectionNone, infSmall
	for _, action := range actions {
		xt, yty := action.Transform()
		block := agent.scope(xt, yty)
		var reward float64
		if block == nil {
			continue
		}
		reward = agent.CarrotWeight * agent.calulateCarrotWeight(xt, yty)
		reward += agent.PythonWeight * agent.calulatePythonWeight(xt, yty)
		reward += agent.ObsticleWeight * agent.calulateObsticleWeight(xt, yty)

		if reward > maxScore {
			maxScore = reward
			d = action
		}
	}
	if maxScore <= infSmall {
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

	agent.prevousAction = d
	return d
}

func (agent *Agent) calulateCarrotWeight(x, y int) float64 {
	_, _, dist := agent.scope.NearestMatching(func(b *pacmound.Block) bool {
		return b != nil && b.Reward() > 0
	}, 10)
	return 1 / dist * dist
}
func (agent *Agent) calulatePythonWeight(x, y int) float64 {
	reward := 0.0
	dist := 2
	for xt := -dist; xt <= dist; xt++ {
		for yt := -dist; yt <= dist; yt++ {
			if !(yt == y && xt == x) {
				b := agent.scope(xt, yt)
				if b != nil && b.IsOccupiedWithPython() && !(yt == y && xt == x) {
					// fmt.Print("*")
					if d := distance(x, y, xt, yt); d > 0 {
						reward += d
					}
				} else {
					// fmt.Print("-")
				}
			} else {
				// fmt.Print("@")
			}
		}
		// fmt.Println()
	}
	if reward < InitalQ {
		reward = InitalQ
	} else {
		reward = 1 / (reward * reward)
	}
	// fmt.Println(1 / (reward * reward))
	return reward
}
func (agent *Agent) calulateObsticleWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		b := agent.scope(x+xt, y+yt)
		if b == nil || b.IsObstructed() || dist > 5 {
			reward += 1
			break
		}
	}
	if reward < InitalQ {
		reward = InitalQ
	} else {
		reward = 1 / (reward * reward)
	}
	// fmt.Println(1 / (reward * reward))
	return reward
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
