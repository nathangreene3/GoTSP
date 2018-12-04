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
	d, p := geneticSoln(importPoints("cities.csv"))
	fmt.Printf("Path: %v\nDist: %0.2f\n", p, d)
}

// geneticSoln solves the TSP by generating a population of permutations,
// then reproducing through a series of generations with a small chance
// to mutate. Reproduction is determined by elitism and PMX is used for
// crossing chromosomes. Mutation occurs by reversing some substring in
// the permutation. The minimum distance and the corresponding
// permutation are returned.
func geneticSoln(ps pointSet) (float64, permutation) {
	pop := randPopulation(10, ps)
	for i := 0; i < 1000000; i++ {
		pop = reproduce(pop, 0.50, 0.25)
	}
	if !isPermutation(pop.perms[0]) {
		log.Fatalf("path not a permution: %v\n", pop.perms[0])
	}
	return totalDist(pop.points, pop.perms[0]), pop.perms[0]
}

// naiveSoln solves the TSP by generating all permuations
// lexicographically and returning the minimum distance and the
// corresponding permutation.
func naiveSoln(ps pointSet) (float64, permutation) {
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
