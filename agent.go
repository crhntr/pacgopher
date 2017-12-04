package pacmound

import (
	"fmt"
	"math"
)

// AgentType represents how the game should handle an agent
// negative agents are pythons
// positive agents are gophers
// zero value is not an angent (empty)
type AgentType int

type AgentGetter func() Agent

type Agent interface {
	Kill()
	Damage(damage Damage)

	CalculateIntent() Direction // player decision loop

	SetScoreGetter(f ScoreGetter)
	SetScopeGetter(f ScopeGetter)
}

type ScoreGetter func() float64
type ScopeGetter func(xOffset, yOffset int) *Block

type AgentData struct {
	a     Agent
	score float64
	x, y  int
	t     AgentType
	dead  bool
}

func (scope ScopeGetter) NearestMatching(match func(b *Block) bool, maxScan int) (minX, minY int, minDist float64) {
	minX, minY = math.MaxInt32/2, math.MaxInt32/2
	minDist = float64(minX * minY)
	check := func(xOffset, yOffset int) {
		if block := scope(xOffset, yOffset); match(block) {
			if dist := distance(0, 0, xOffset, yOffset); dist < minDist {
				minX, minY, minDist = xOffset, yOffset, dist
			}
		}
	}
	for dist := 1; dist <= maxScan; dist++ {
		for i := -dist; i < dist; i++ {
			check(i, -dist)
			check(-dist, i)
			check(i, dist)
			check(dist, i)
		}
	}
	return
}

func (scope ScopeGetter) DisplayRegion(dist int) {
	fmt.Println()
	for x := -dist; x <= dist; x++ {
		for y := -dist; y <= dist; y++ {
			b := scope(x, y)
			if b == nil {
				fmt.Print("    ")
			} else if x == 0 && y == 0 {
				fmt.Print("[@@]")
			} else if b.IsOccupied() {
				fmt.Print("[PY]")
			} else if b.IsObstructed() {
				fmt.Print("[##]")
			} else if r := b.Reward(); r > 0 {
				fmt.Printf("[$$]")
			} else {
				fmt.Printf("[--]")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (ad *AgentData) IsPython() bool {
	return ad.t < 0
}

func (ad *AgentData) IsGopher() bool {
	return ad.t > 0
}

func (ad *AgentData) Score() float64 {
	return ad.score
}

func newScopeGetter(maze Maze, ad *AgentData) ScopeGetter {
	return func(xOffset, yOffset int) *Block {

		x := ad.x + xOffset
		y := ad.y + yOffset

		if x < 0 || y < 0 || x >= len(maze) || y >= len(maze[x]) {
			return nil
		}

		return &Block{
			agent:    maze[x][y].agent,
			obsticle: maze[x][y].obsticle,
			reward:   maze[x][y].reward,
		}
	}
}

func (maze *Maze) setAgent(x, y int, agent Agent) (*AgentData, error) {
	if x < 0 || y < 0 || x >= len(*maze) || y >= len((*maze)[x]) {
		return &AgentData{}, ErrBeyondTheKnownMaze
	}
	(*maze)[x][y].agent = &AgentData{
		a:     agent,
		x:     x,
		y:     y,
		score: 0,
	}
	return (*maze)[x][y].agent, nil
}

func fight(agent1, agent2 *AgentData) bool {
	if agent1.dead || agent2.dead {
		return true
	}
	if agent1.IsGopher() && agent2.IsPython() {
		agent1.Damage(DamageLostFight)
		agent2.score += agent1.score
		return true
	}
	if agent2.IsGopher() && agent1.IsPython() {
		agent2.Damage(DamageLostFight)
		agent1.score += agent2.score
		return true
	}
	return false
}
