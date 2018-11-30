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
	d, p := geneticSoln(getPoints("fourpoints.csv"))
	fmt.Printf("Path: %v\nDist: %0.2f\n", p, d)
}

func geneticSoln(ps pointSet) (float64, permutation) {
	pop := randPopulation(4, ps)
	for i := 0; i < 10; i++ {
		for j := range pop.perms {
			if !isPerm(pop.perms[j]) {
				fmt.Printf("gen %d, index %d, path not a permution: %v\n", i, j, pop.perms[j])
			}
		}
		pop = reproduce(pop, 0.25, 0.10)
	}
	if !isPerm(pop.perms[0]) {
		log.Fatalf("path not a permution%v\n", pop)
	}
	return totalDist(pop.points, pop.perms[0]), pop.perms[0]
}

// naiveSoln solves the TSP by generating all permuations
// lexicographically and returning the minimum distance and the
// corresponding permutation.
func naiveSoln(ps pointSet) (float64, permutation) {
	perm := basePerm(len(ps))  // Current permutation to try
	base := copyPerm(perm)     // Used to check if all permutations have been tried
	minPerm := copyPerm(perm)  // Permutation resulting in minimum distance
	dist := 0.0                // Current distance
	minDist := math.MaxFloat64 // Minimum distance found

	// Try all n! permutations and store the current best solution
	for {
		dist = totalSqDist(ps, perm)
		if dist < minDist {
			minDist = dist
			copy(minPerm, perm)
		}
		perm = nextPerm(perm)
		if comparePerms(perm, base) {
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
