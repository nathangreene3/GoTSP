package main

import (
	"math/rand"
	"sort"
	"time"
)

// permutation is an ordering on Z(n).
type permutation []int

// population is a set of permutations.
type population struct {
	perms  []permutation
	points pointSet
}

// nextPerm returns the next lexicographical permutation. If
// the current permutation is the last permutation, then the
// base permutation is returned.
func nextPerm(p permutation) permutation {
	// Find largest index k such that p[k] < p[k+1]. If k
	// remains -1, p is the final lexicographical
	// permutation (ie n-1...210).
	n := len(p)
	k := -1
	for i := n - 2; 0 <= i; i-- {
		if p[i] < p[i+1] {
			k = i
			break
		}
	}

	// Return the first permuation if no k found.
	if k == -1 {
		return basePerm(n)
	}

	// Find largest index j > k such that p[k] < p[j].
	j := -1
	for i := n - 1; k < i; i-- { // 0 <= k < n-1
		if p[k] < p[i] {
			j = i
			break
		}
	}

	// Swap p[k] and p[j].
	q := copyPerm(p)
	q[k], q[j] = q[j], q[k]

	// Reverse p[k+1:].
	a, b := k+1, n-1
	for a < b {
		q[a], q[b] = q[b], q[a]
		a++
		b--
	}

	return q
}

// basePerm returns the base permutation (012...n-1).
func basePerm(n int) permutation {
	p := make(permutation, n)
	for i := range p {
		p[i] = i
	}
	return p
}

// randPerm returns a random permutation of length n.
func randPerm(n int) permutation {
	rand.Seed(int64(time.Now().Second()))
	p := basePerm(n)
	var a, b int
	for i := 0; i < 3*n; i++ {
		a, b = rand.Intn(n), rand.Intn(n)
		p[a], p[b] = p[b], p[a]
	}
	return p
}

// heapPermute is taken from Analysis of Algorithms, page 54.
func heapPermute(a permutation, n int) permutation {
	// TODO: make this functional.
	if 0 < n {
		for i := 0; i < n; i++ {
			a = heapPermute(a, n-1)
			if n%2 == 0 {
				a[i], a[n-1] = a[n-1], a[i]
			} else {
				a[0], a[n-1] = a[n-1], a[0]
			}
		}
	}
	return a
}

// copyPerm returns a new permution copy.
func copyPerm(x permutation) permutation {
	y := make(permutation, len(x))
	copy(y, x)
	return y
}

// compareIntSlices returns true if the two slices are nil
// or if they share the same length and values at each index.
func comparePerms(x, y permutation) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// **********************************************************************
// Genetic algorithm functions
// **********************************************************************

// cross ...TODO
func cross(x, y permutation) (permutation, permutation) {
	n := len(x)
	u, v := make(permutation, n), make(permutation, n)

	// Method 1: swap ends on pivot
	pivot := rand.Intn(n)
	copy(u[:pivot], x[:pivot])
	copy(v[:pivot], y[:pivot])
	copy(u[pivot:], y[pivot:])
	copy(v[pivot:], x[pivot:])

	// Method 2: swap middle on start, end
	// end := rand.Intn(n-1) + 1   // 0 < end < n
	// start := rand.Intn(end + 1) // 0 <= start <= end
	// copy(u[:start], x[:start])
	// copy(u[start:end], y[start:end])
	// copy(u[end:], x[end:])
	// copy(v[:start], y[:start])
	// copy(v[start:end], x[start:end])
	// copy(v[end:], y[end:])

	return u, v
}

// mutate ...TODO
func mutate(p permutation) permutation {
	n := len(p)
	b := rand.Intn(n-1) + 1 // 0 < b < n
	a := rand.Intn(b)       // 0 <= a < b
	q := copyPerm(p)

	// Reverse subsequence in permutation
	for a < b {
		q[a], q[b] = q[b], q[a]
		a++
		b--
	}
	return q
}

// reproduce ...TODO
func reproduce(pop *population) *population {
	sort.Sort(pop)
	return pop
}

// randPopulation ...TODO
func randPopulation(size int, ps pointSet) *population {
	n := len(ps)
	pop := &population{
		perms: make([]permutation, n),
	}
	copy(pop.points, ps)
	for i := range pop.perms {
		pop.perms[i] = randPerm(n)
	}
	return pop
}

func (pop *population) Len() int {
	return len(pop.perms)
}

func (pop *population) Less(i, j int) bool {
	return totalDist(pop.points, pop.perms[i]) < totalDist(pop.points, pop.perms[j])
}

func (pop *population) Swap(i, j int) {
	pop.perms[i], pop.perms[j] = pop.perms[j], pop.perms[i]
}
