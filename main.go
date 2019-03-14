package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().UnixNano()))

	ps, err := importPoints("cities.csv")
	if err != nil {
		log.Fatal(err)
	}

	dist, perm := geneticSoln(ps, 10, 1000000, reverseSubsequence, pmx)
	fmt.Printf("Path: %v\nDist: %0.2f\n", perm, dist)
}

func run() error {
	// flg := flag.NewFlagSet("GoTSP", flag.ContinueOnError)
	// input := flg.String("input_file", "", "csv file containing points")
	// shortestPerm := flg.String("shortest_path_file", "", "csv file containing shortest known permutation through input file")
	// generations := flg.Int("generations", 1000000, "number of generations")
	// popSize := flg.Int("population_size", 10, "number of members in a population")

	// err := flg.Parse(os.Args[1:])
	// if err != nil {
	// 	return err
	// }

	return nil
}

// geneticSoln solves the TSP by generating a population of permutations, then
// reproducing through a series of generations with a small chance to mutate.
// Reproduction is determined by elitism. The minimum distance found and the
// corresponding permutation are returned.
func geneticSoln(ps pointSet, popSize int, generations int, f mutateFunc, g breedFunc) (float64, permutation) {
	currentShortestPath, err := importPermutation("shortestpath.csv")
	if err != nil {
		log.Fatal(err)
	}

	pop := randPopulation(popSize, ps)
	if currentShortestPath != nil && len(currentShortestPath) == len(ps) {
		copy(pop.perms[0], currentShortestPath)
	}
	sort.Sort(pop)

	for i := 0; i < generations; i++ {
		pop = reproduce(pop, 0.50, 0.25, f, g)
	}
	for i, p := range pop.perms {
		fmt.Println(i, p)
	}

	if !isPermutation(pop.shortestPerm) {
		log.Fatalf("path not a permution: %v\n", pop.shortestPerm)
	}

	err = pop.shortestPerm.exportPermutation("shortestpath.csv")
	if err != nil {
		log.Fatal(err)
	}

	return totalDist(pop.points, pop.shortestPerm), pop.shortestPerm
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
		dist = totalDist(ps, perm)
		if dist < minDist {
			minDist = dist
			minPerm = copyPermutation(perm)
		}

		perm = nextPermutation(perm)
		if isBase(perm) {
			break
		}
	}

	return minDist, minPerm
}

// factorial returns n! for all n > 1 and 1 for all n < 2.
func factorial(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}

	return f
}
