package pacmound

import (
	"fmt"
	"math/rand"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func prob(p float64) bool {
	return rand.Float64() <= p
}

func displayQS(qs [][][]float64) {
	for i := range qs {
		for j := range qs[i] {
			fmt.Print("[ ")
			for _, val := range qs[i][j] {
				fmt.Printf("%4.01f ", val)
			}
			fmt.Print("]  ")
		}
		fmt.Println()
	}
}

func sum(qs [][][]float64) float64 {
	sum := 0.0
	for i := range qs {
		for j := range qs[i] {
			for k := range qs[i][j] {
				sum += qs[i][j][k]
			}
		}
	}
	return sum
}

func copyValues(dst, src [][][]float64) {
	for i := range src {
		for j := range src[i] {
			for k := range src[i][j] {
				dst[i][j][k] = src[i][j][k]
			}
		}
	}
}
