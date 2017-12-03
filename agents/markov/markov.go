package markov

import (
	"math"
	"math/rand"

	"github.com/crhntr/pacmound"
	"github.com/crhntr/pacmound/agents"
)

const (
	numberOfActions = 4 + 1
	InitalQ         = 0.00001
)

var (
	infSmall = math.Inf(-1)
)

type Agent struct {
	DiscountFactor,
	LearningRate float64

	check  pacmound.ScopeGetter
	reward pacmound.ScoreGetter

	prevousAction pacmound.Direction

	QTable QTable

	previousRewardTotal,
	carrotWeight,
	obsticalWeight,
	pythonWeight float64
}

func (agent *Agent) SetScoreGetter(f pacmound.ScoreGetter) { agent.reward = f }
func (agent *Agent) SetScopeGetter(f pacmound.ScopeGetter) { agent.check = f }
func (agent *Agent) Warning(err error)                     {}
func (agent *Agent) CalculateIntent() pacmound.Direction {
	return agent.QLearning()
}

func (agent *Agent) QLearning() pacmound.Direction {
	if agent.carrotWeight == 0 || agent.pythonWeight == 0 {
		return pacmound.Direction(rand.Intn(4) + 1)
	}

	rewardCurrent := agent.reward() - agent.previousRewardTotal
	agent.previousRewardTotal = agent.reward()

	x, y := agent.prevousAction.Transform()
	agent.QTable.update(rewardCurrent, agent.LearningRate, agent.DiscountFactor, x*-1, y*-1, agent.prevousAction)

	actions := agents.Actions()
	for xOffset := 0 - (len(agent.QTable) / 2); xOffset < len(agent.QTable); xOffset++ {
		for yOffset := 0 - (len(agent.QTable[xOffset]) / 2); yOffset < len(agent.QTable[xOffset]); yOffset++ {
			for _, action := range actions[1:] {
				xt, yt := action.Transform()
				block := agent.check(x, y)

				reward := infSmall
				if block != nil {
					reward = agent.carrotWeight * agent.calulateCarrotWeight(xOffset*xt, yOffset*yt)
					reward += agent.pythonWeight * agent.calulatePythonWeight(xOffset*xt, yOffset*yt)
					reward += agent.obsticalWeight * agent.calulateObsticalWeight(xOffset*xt, yOffset*yt)
				}

				agent.QTable.update(reward, agent.LearningRate, agent.DiscountFactor, xOffset, yOffset, action)
			}
		}
	}

	q := agent.carrotWeight * agent.calulateCarrotWeight(0, 0)
	q += agent.pythonWeight * agent.calulatePythonWeight(0, 0)
	q += agent.obsticalWeight * agent.calulateObsticalWeight(0, 0)

	agent.carrotWeight = agent.carrotWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulateCarrotWeight(x, y)
	agent.pythonWeight = agent.pythonWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulatePythonWeight(x, y)
	agent.obsticalWeight = agent.obsticalWeight + agent.LearningRate*(rewardCurrent-q)*agent.calulateObsticalWeight(x, y)

	agent.prevousAction = agent.QTable[0][0].maxDirection()
	return agent.prevousAction
}

func (q QTable) update(r, α, γ float64, x, y int, action pacmound.Direction) {
	q.setQ(x, y, action, (1-α)*q.getQ(x, y, action)+α*(r+γ*q[x][y].maxReward()))
}

func (agent *Agent) calulateCarrotWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		for {
			b := agent.check(x+xt*dist, y+yt*dist)
			if b == nil || b.IsObstructed() {
				break
			}
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
			if b == nil || b.IsObstructed() {
				break
			}
			if b.IsOccupied() {
				reward += 1 / float64(dist*dist)
			}
		}
	}
	return reward
}
func (agent *Agent) calulateObsticalWeight(x, y int) float64 {
	reward := 0.0
	actions := agents.Actions()
	for _, action := range actions[1:] {
		dist := 1
		xt, yt := action.Transform()
		for {
			b := agent.check(x+xt*dist, y+yt*dist)
			if b == nil || b.IsObstructed() {
				reward += 1 / float64(dist*dist)
				break
			}
		}
	}
	return reward
}
