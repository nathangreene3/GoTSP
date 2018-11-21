package main

import (
	"fmt"
	"math"
)

func main() {
	ps := pointSet{
		point{0, 0},
		point{2, 2},
		point{3, 1},
		point{4, 2},
	}
	d, p := naiveSoln(ps)
	fmt.Printf("Path: %v, Dist: %0.2f\n", p, d)
	// n := 5
	// p := basePerm(n)
	// for i := 0; i < factorial(n); i++ {
	// 	fmt.Printf("%v\n", p)
	// 	p = nextPerm(p)
	// }
}

// factorial returns n!
func factorial(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	return f
}

func naiveSoln(ps pointSet) (float64, permutation) {
	perm := basePerm(len(ps))  // Current permutation to try
	base := copyPerm(perm)     // Used to check if all permutations have been tried
	minPerm := copyPerm(perm)  // Permutation resulting in minimum distance
	dist := 0.0                // Current distance
	minDist := math.MaxFloat64 // Minimum distance found
	for {
		dist = totalSqDist(ps, perm)
		fmt.Printf("%v %0.2f %v %0.2f\n", perm, dist, minPerm, minDist)
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
