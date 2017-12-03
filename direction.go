package pacmound

import (
	"math/rand"
	"strings"
)

type Direction int

const (
	DirectionNone Direction = iota
	DirectionNorth
	DirectionEast
	DirectionSouth
	DirectionWest
)

var dirVector = [][]int{{0, 0},
	{1, 0}, {0, 1},
	{-1, 0}, {0, -1},
}

func (dir Direction) TurnLeft() Direction {
	return Direction((int(dir) + (len(dirVector) - 1)) % len(dirVector))
}

func (dir Direction) TurnRight() Direction {
	return Direction((int(dir) + 1) % len(dirVector))
}

func RandomDirection() Direction {
	return Direction(rand.Intn(len(dirVector)))
}

func (dir Direction) String() string {
	return []string{"?", "v", ">", "^", "<"}[dir]
}

func ParseDirection(str string) Direction {
	index := strings.Index("?v>^<", str)
	if index < 0 {
		return DirectionNone
	}
	return directions[index]
}

func (dir Direction) Transform() (x, y int) {
	return dirVector[dir][0], dirVector[dir][1]
}

func move(m *Maze, dir Direction, x, y int) (int, int) {
	if dir < 0 {
		return x, y
	}
	nextX, nextY := x+dirVector[dir][0], y+dirVector[dir][1]
	if nextX < 0 {
		nextX = 0
	}
	if nextY < 0 {
		nextY = 0
	}
	if nextX >= len((*m)) {
		nextX = nextX % len((*m))
	}
	if nextY >= len((*m)[nextX]) {
		nextY = nextY % len((*m)[nextX])
	}
	return nextX, nextY
}

var directions = []Direction{
	DirectionNorth,
	DirectionEast,
	DirectionSouth,
	DirectionWest,
}
