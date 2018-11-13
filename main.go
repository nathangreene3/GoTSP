package main

import (
	"fmt"
)

func main() {
	p := basePerm(4)
	fmt.Println(nextPerm(p))
	// fmt.Println(*p)
	// for i := 0; i < factorial(len(*p)); i++ {
	// 	*p = nextPerm(*p)
	// 	fmt.Println(*p)
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

func removeCrossPaths(ps pointSet, perm permutation) {

}

func naiveSoln(ps pointSet) (minDist float64, minPerm permutation) {
	// perm := basePerm(len(ps))
	// base := copyPerm(perm)
	// minPerm = copyPerm(perm)
	// dist := 0.0
	// minDist = math.MaxFloat64
	// notEqual := true
	// for notEqual {
	// 	dist = totalSqDist(ps, perm)
	// 	if dist < minDist {
	// 		minDist = dist
	// 		minPerm = copyPerm(perm)
	// 	}
	// 	perm = nextPerm(perm)
	// 	if comparePerms(perm, base) {
	// 		notEqual = false
	// 	}
	// }
	// minDist = totalDist(ps, minPerm)
	return minDist, minPerm
}
