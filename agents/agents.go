package agents

import "github.com/crhntr/pacmound"

func Actions() []pacmound.Direction {
	return []pacmound.Direction{
		pacmound.DirectionNorth,
		pacmound.DirectionEast,
		pacmound.DirectionSouth,
		pacmound.DirectionWest,
	}
}

func removeMinimumScoringDirections(scores []float64, directions []pacmound.Direction) ([]float64, []pacmound.Direction) {
	if allEqual(scores) {
		return scores, directions
	}
	max, minIndex := scores[0], 0
	i := 0
	for ; i < len(scores); i++ {
		if scores[i] < max {
			max = scores[i]
			minIndex = i
		}
	}
	scores = append(scores[:minIndex], scores[minIndex+1:]...)
	directions = append(directions[:minIndex], directions[minIndex+1:]...)
	return removeMinimumScoringDirections(scores, directions)
}

func allEqual(scores []float64) bool {
	for i := 1; i < len(scores); i++ {
		if scores[i] != scores[i-1] {
			return false
		}
	}
	return true
}
