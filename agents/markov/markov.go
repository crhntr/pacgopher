package markov

import (
	"math/rand"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
)

const (
	numberOfActions = 4 + 1
	InitalQ         = 0.00001
)

var (
	infSmall = -100000000.0
)

type Agent struct {
	DiscountFactor float64 `json:"discountFactor"`
	LearningRate   float64 `json:"learningRate"`

	CarrotWeight   float64 `json:"carrotWeight"`
	ObsticleWeight float64 `json:"obsticleWeight"`
	PythonWeight   float64 `json:"pythonWeight"`

	check  pacmound.ScopeGetter
	reward pacmound.ScoreGetter

	prevousAction pacmound.Direction

	previousRewardTotal float64

	count int
}

func (agent *Agent) SetScoreGetter(f pacmound.ScoreGetter) { agent.reward = f }
func (agent *Agent) SetScopeGetter(f pacmound.ScopeGetter) { agent.check = f }
func (agent *Agent) Warning(err error)                     {}
func (agent *Agent) CalculateIntent() pacmound.Direction {
	defer func() { agent.count++ }()
	if agent.count < 10 {
		return pacmound.Direction(rand.Intn(4) + 1)
	}
	return agent.QLearning(agent.LearningRate, agent.DiscountFactor)
}
func (agent *Agent) Kill() {
	agent.QLearning(agent.LearningRate, agent.DiscountFactor)
}

func (agent *Agent) QLearning(α, γ float64) pacmound.Direction {
	rewardCurrent := agent.reward() - agent.previousRewardTotal
	agent.previousRewardTotal = agent.reward()

	xt, yt := agent.prevousAction.Transform()

	actions := agents.Actions()[1:]
	rewards := make([]float64, len(actions))
	for i, action := range actions[1:] {
		xtt, ytt := action.Transform()
		block := agent.check(xtt, ytt)
		reward := infSmall
		if block != nil {
			reward = agent.CarrotWeight * agent.calulateCarrotWeight(xt*xtt, yt*ytt)
			reward += agent.PythonWeight * agent.calulatePythonWeight(xt*xtt, yt*ytt)
			reward += agent.ObsticleWeight * agent.calulateObsticleWeight(xt*xtt, yt*ytt)
		}
		rewards[i] = reward
	}

	q := agent.CarrotWeight * agent.calulateCarrotWeight(0, 0)
	q += agent.PythonWeight * agent.calulatePythonWeight(0, 0)
	q += agent.ObsticleWeight * agent.calulateObsticleWeight(0, 0)

	agent.CarrotWeight = agent.CarrotWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulateCarrotWeight(xt, yt)
	agent.PythonWeight = agent.PythonWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulatePythonWeight(xt, yt)
	agent.ObsticleWeight = agent.ObsticleWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulateObsticleWeight(xt, yt)

	agent.prevousAction = maxDirection(rewards...) + 1
	return agent.prevousAction
}

func (agent *Agent) calulateCarrotWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		for {
			b := agent.check(x+xt*dist, y+yt*dist)
			if b == nil || b.IsObstructed() || dist > 5 {
				break
			}
			dist++
			reward += b.Reward() / float64(dist*dist)
		}
	}
	return reward
}
func (agent *Agent) calulatePythonWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		for {
			b := agent.check(x+xt*dist, y+yt*dist)
			if b == nil || b.IsObstructed() || dist > 5 {
				break
			}
			dist++
			if b.IsOccupied() {
				reward += 1 / float64(dist)
			}
		}
	}
	return reward
}
func (agent *Agent) calulateObsticleWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		for {
			b := agent.check(x+xt*dist, y+yt*dist)
			if b == nil || b.IsObstructed() || dist > 5 {
				reward += 1 / float64(dist)
				break
			}
			dist++
		}
	}
	return reward
}

func maxReward(rewards ...float64) float64 {
	maxReward := rewards[0]
	for i := 1; i < len(rewards); i++ {
		if rewards[i] > maxReward {
			maxReward = rewards[i]
		}
	}
	return maxReward
}

func maxDirection(rewards ...float64) pacmound.Direction {
	maxDir, maxReward := 0, rewards[0]
	for i := 1; i < len(rewards); i++ {
		if rewards[i] > maxReward {
			maxDir, maxReward = i, rewards[i]
		}
	}
	return pacmound.Direction(maxDir)
}
