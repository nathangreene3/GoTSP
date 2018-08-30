package main

import (
	"math"
)

type vector struct {
	id     int
	values []float64
}

type path []vector
type permutation []int

func main() {
	n := 4
	p := make(path, n)
	for i := 0; i < n; i++ {
		p[i].id = i
	}
	p[0].values = []float64{0, 0}
	p[1].values = []float64{2, 2}
	p[2].values = []float64{3, 1}
	p[3].values = []float64{4, 2}
	solve(p)
}

func solve(p path) {
	// n := len(p)
	// perm := nextPerm(make(permutation, n))
	// minDist := totalDistance(p)
	// dist:=minDist
	// for i:=0;i<fact(n);i++{

	// }

	// fmt.Println(minDist)
}

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

func totalDistance(p path) (d float64) {
	n := len(p)
	for i := 0; i+1 < n; i++ {
		d += math.Sqrt(sqDist(p[i], p[i+1]))
	}
	d += math.Sqrt(sqDist(p[0], p[n-1]))
	return d
}

func totalSqDist(p path) (d float64) {
	n := len(p)
	for i := 0; i+1 < n; i++ {
		d += sqDist(p[i], p[i+1])
	}
	d += sqDist(p[0], p[n-1])
	return d
}

func sqDist(U, V vector) (d float64) {
	n := len(U.values)
	for i := 0; i < n; i++ {
		d += (U.values[i] - V.values[i]) * (U.values[i] - V.values[i])
	}
	return d
}

func fact(n int) (f int) {
	f = 1
	for i := 2; i < n; i++ {
		f *= i
	}
	f *= n
	return f
}
