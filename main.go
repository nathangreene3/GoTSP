package main

import (
	"math"
)

type point []float64
type pointSet []point
type permutation []int

func main() {

}

// nextPerm returns the next lexicographical permutation. If
// the current permutation is the last permutation, then the
// base permutation is returned.
func nextPerm(p permutation) (q permutation) {
	// 1. Find largest index k such that p[k] < p[k+1]. If k
	//    remains -1, p is the final lexicographical
	//    permutation (ie n-1...210).
	n := len(p)
	q = make(permutation, n)
	k := -1
	for i := 0; i+1 < n; i++ {
		if p[i] < p[i+1] {
			k = i
		}
	}

	// 2.a Find largest index j > k such that p[k] < p[j]. If
	//     k = -1, return the first permutation (012...n-1).
	j := -1
	if 0 <= k {
		for i := k + 1; i < n; i++ {
			if p[k] < p[i] {
				j = i
			}
		}

		// 3. Swap p[k] and p[j].
		for i := 0; i < k; i++ {
			q[i] = p[i]
		}
		for i := k + 1; i < n; i++ {
			q[i] = p[n-i-1]
		}
		q[k] = p[j]
		q[j] = p[k]
	} else {
		// 2.b Generate the first permuation.
		for i := 0; i < len(p); i++ {
			q[i] = i
		}
	}
	return q
}

// totalDist returns the total distance traversed across a
// path of points. A path is a permutation of points.
// Assumes the point set and the permutation are of equal
// dimension.
func totalDist(ps pointSet, perm permutation) (d float64) {
	n := len(ps)
	for i := 0; i+1 < n; i++ {
		d += math.Sqrt(sqDist(ps[perm[i]], ps[perm[i+1]]))
	}
	d += math.Sqrt(sqDist(ps[perm[0]], ps[perm[n-1]]))
	return d
}

// totalSqDist returns the total squared distance traversed
// across a path of points. A path is a permutation of
// points. Assumes teh point set and the permutation are of
// equal dimension.
func totalSqDist(ps pointSet, perm permutation) (d float64) {
	n := len(ps)
	for i := 0; i+1 < n; i++ {
		d += sqDist(ps[perm[i]], ps[perm[i+1]])
	}
	d += sqDist(ps[perm[0]], ps[perm[n-1]])
	return d
}

// sqDist returns the squared distance between two points.
// Assumes the points are equal in dimension.
func sqDist(x, y point) (d float64) {
	for i := range x {
		d += (x[i] - y[i]) * (x[i] - y[i])
	}
	return d
}
