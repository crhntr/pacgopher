package markov

import "github.com/crhntr/pacmound"

type ActionRewards [numberOfActions]float64
type QTable [][]ActionRewards

func (gt QTable) Empty() bool {
	return len(gt) == 0
}

func NewQTable(xSize, ySize int) QTable {
	table := make([][]ActionRewards, xSize)
	for x := 0; x < ySize; x++ {
		table[x] = make([]ActionRewards, ySize)
		for y := 0; y < ySize; y++ {
			for d := range table[x][y] {
				table[x][y][d] = InitalQ
			}
		}
	}
	return QTable(table)
}

func (gt QTable) getQ(x, y int, dir pacmound.Direction) float64 {
	x += len(gt) / 2
	y += len(gt[x]) / 2
	if x < 0 || x >= len(gt) || y < 0 || y >= len(gt[x]) {
		return infSmall
	}
	return gt[x][y][dir]
}

func (gt QTable) setQ(x, y int, dir pacmound.Direction, newQ float64) {
	x += len(gt) / 2
	y += len(gt[x]) / 2

	gt[x][y][dir] = newQ
}

func (ar ActionRewards) maxReward() float64 {
	maxReward := ar[0]
	for i := 1; i < len(ar); i++ {
		if ar[i] > maxReward {
			maxReward = ar[i]
		}
	}
	return maxReward
}

func (ar ActionRewards) maxDirection() pacmound.Direction {
	maxDir, maxReward := 0, ar[0]
	for i := 1; i < len(ar); i++ {
		if ar[i] > maxReward {
			maxDir, maxReward = i, ar[i]
		}
	}
	return pacmound.Direction(maxDir)
}
