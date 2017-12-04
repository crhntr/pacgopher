package pacmound

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
)

var (
	ErrBeyondTheKnownMaze = errors.New(
		"beyond the known maze")
	ErrCrashedIntoObsticle = errors.New(
		"unessary headache, crashed into obstacle")
	ErrIncompletePolicyDeclaration = errors.New(
		"incomplete policy declaration incorrect number of directions passed")
)

// Maze should be initalized with the NewMaze function
type Maze [][]Block

func (m Maze) encodable() [][]EncodedBlock {
	eb := make([][]EncodedBlock, len(m))
	for x := range m {
		eb[x] = make([]EncodedBlock, len(m[x]))
		for y := range m[x] {
			eb[x][y] = m[x][y].encodable()
		}
	}
	return eb
}

func (m Maze) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.encodable())
}

func NewEmptyMaze(xLen, yLen int) Maze {
	if xLen < 0 || yLen < 0 {
		panic("world must be at least 1 by 1")
	}
	maze := make([][]Block, xLen)
	for i := range maze {
		maze[i] = make([]Block, yLen)
	}
	return Maze(maze)
}

func (maze *Maze) setObsticle(x, y int) error {
	if x < 0 || y < 0 || x >= len(*maze) || y >= len((*maze)[x]) {
		return ErrBeyondTheKnownMaze
	}
	(*maze)[x][y].obsticle = true
	(*maze)[x][y].reward = 0
	return nil
}

func (maze *Maze) setReward(x, y int, reward float64) error {
	if x < 0 || y < 0 || x >= len(*maze) || y >= len((*maze)[x]) {
		return ErrBeyondTheKnownMaze
	}
	(*maze)[x][y].obsticle = false
	(*maze)[x][y].reward = reward
	return nil
}

func (maze *Maze) getReward(x, y int) float64 {
	reward := (*maze)[x][y].reward
	(*maze)[x][y].reward = 0
	return reward
}

// RewardAt may returns an error for out of range or obsticles
func (maze Maze) RewardAt(x, y int) (float64, error) {
	if x < 0 || y < 0 || x >= len(maze) || y >= len(maze[x]) {
		return 0, ErrBeyondTheKnownMaze
	}
	if maze[x][y].obsticle {
		return float64(DamageColision), ErrCrashedIntoObsticle
	}
	return maze[x][y].reward, nil
}

// IsObsticle returns if location has an obstical
func (maze *Maze) IsObsticle(x, y int) bool {
	if x < 0 || y < 0 || x >= len(*maze) || y >= len((*maze)[x]) {
		return true // if out of range, we can't get there
	}
	return (*maze)[x][y].obsticle
}

// IsObsticle returns if location has an obstical
func (maze *Maze) Occupant(x, y int) *AgentData {
	if x < 0 || y < 0 || x >= len(*maze) || y >= len((*maze)[x]) || (*maze)[x][y].agent == nil {
		return nil
	}
	return (*maze)[x][y].agent
}

// RemainingReward returns if location has an obstical
func (m Maze) RemainingReward() float64 {
	reward := 0.0
	for x := range m {
		for y := range m[x] {
			reward += m[x][y].reward
		}
	}
	return reward
}

func (maze *Maze) RandomEmptyPosition() (x, y int) {
	attempts := 0
	for {
		x := rand.Intn(len((*maze)))
		y := rand.Intn(len((*maze)[x]))
		if !(*maze)[x][y].obsticle {
			return x, y
		}
		if attempts > len(*maze)*len((*maze)[0]) {
			panic("no empty space")
		}
		attempts++
	}
}

func distance(xFinal, yFinal, xInital, yInital int) float64 {
	return math.Sqrt(float64((xInital-xFinal)*(xInital-xFinal) + (yInital-yFinal)*(yInital-yFinal)))
}
