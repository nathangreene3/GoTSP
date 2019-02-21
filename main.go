package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().UnixNano()))

	ps, err := importPoints("cities.csv")
	if err != nil {
		log.Fatal(err)
	}
	d, p := geneticSoln(ps, 10, 1000000, mutate, cross)
	fmt.Printf("Path: %v\nDist: %0.2f\n", p, d)
}

// geneticSoln solves the TSP by generating a population of permutations, then
// reproducing through a series of generations with a small chance to mutate.
// Reproduction is determined by elitism. The minimum distance found and the
// corresponding permutation are returned.
func geneticSoln(ps pointSet, popSize int, generations int, f mutateFunc, g breedFunc) (float64, permutation) {
	pop := randPopulation(popSize, ps)
	for i := 0; i < generations; i++ {
		pop = reproduce(pop, 0.50, 0.25, f, g)
	}

	p := pop.perms[0]
	if !isPermutation(p) {
		log.Fatalf("path not a permution: %v\n", p)
	}

	err := p.exportPermutation("bestsolution.csv")
	if err != nil {
		log.Fatal(err)
	}

	return totalDist(pop.points, p), p
}

// naiveSoln solves the TSP by generating all permuations lexicographically and
// returning the minimum distance and the corresponding permutation. This will
// fail if the point set contains more than 15 points.
func naiveSoln(ps pointSet) (float64, permutation) {
	if 15 < len(ps) {
		log.Fatalf("cannot solve on more than 15 points in a reasonable amount of time")
	}

	perm := basePermutation(len(ps)) // Current permutation to try
	minPerm := copyPermutation(perm) // Permutation resulting in minimum distance
	dist := 0.0                      // Current distance
	minDist := math.MaxFloat64       // Minimum distance found

	// Try all n! permutations and store the current best solution
	for {
		dist = totalSqDist(ps, perm)
		if dist < minDist {
			minDist = dist
			minPerm = copyPermutation(perm)
		}

		perm = nextPermutation(perm)
		if isBase(perm) {
			break
		}
	}

	return totalDist(ps, minPerm), minPerm
}

// factorial returns n!
func factorial(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}

	return f
}
