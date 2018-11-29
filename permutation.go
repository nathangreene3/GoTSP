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

// copyPerm returns a new permution copy.
func copyPerm(x permutation) permutation {
	y := make(permutation, len(x))
	copy(y, x)
	return y
}

// copyPerm returns a new permutation copy.
func (p permutation) copyPerm() permutation {
	q := make(permutation, len(p))
	copy(q, p)
	return q
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

func isPerm(a []int) bool {
	sort.Ints(a)
	for i := range a {
		if i != a[i] {
			return false
		}
	}
	return true
}

/*
*************************************************************************
	Genetic algorithm functions
*************************************************************************
*/

// cross returns two new permutions each leading with the values of x and
// y but with trailing values of y and x. The pivot point is selected at
// random. For example, [1,2,3,4] and [5,6,7,8] might be crossed at index
// 1 returning [1,6,7,8] and [5,2,3,4].
func cross(x, y permutation) (permutation, permutation) {
	n := len(x)
	u, v := make(permutation, n), make(permutation, n)

	// Swap ends on pivot
	rand.Seed(int64(time.Now().Second()))
	pivot := rand.Intn(n)
	copy(u[:pivot], x[:pivot])
	copy(v[:pivot], y[:pivot])
	copy(u[pivot:], y[pivot:])
	copy(v[pivot:], x[pivot:])

	return u, v
}

// mutate reverses a subsequence within a permutation. For example,
// [1,2,3,4] might be mutated as [1,3,2,4], or even [3,2,1,4].
func mutate(p permutation) permutation {
	b := rand.Intn(len(p)-1) + 1 // 0 < b < n
	a := rand.Intn(b)            // 0 <= a < b
	q := copyPerm(p)

	// Reverse subsequence in permutation
	for a < b {
		q[a], q[b] = q[b], q[a]
		a++
		b--
	}

	return q
}

// reproduce
func reproduce(pop *population, topPct, mutationRate float64) *population {
	rand.Seed(int64(time.Now().Second()))
	nextGen := copyPopulation(pop)
	n := len(nextGen.perms)
	for i := 0; i < int(topPct*float64(n)); i += 2 {
		nextGen.perms[n-i-1], nextGen.perms[n-i-2] = cross(nextGen.perms[i], nextGen.perms[i+1])
		if rand.Float64() < mutationRate {
			nextGen.perms[n-i-1], nextGen.perms[n-i-2] = mutate(nextGen.perms[n-i-1]), mutate(nextGen.perms[n-i-2])
		}
	}
	sort.Sort(nextGen)

	return nextGen
}

// randPopulation returns a random population of a given size over a set
// of points.
func randPopulation(size int, ps pointSet) *population {
	if size <= 0 {
		panic("population size must be positive")
	}

	pop := &population{
		perms:  make([]permutation, size),
		points: copyPointSet(ps),
	}

	n := len(ps)
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

// copyPopulation returns a copy of a given population.
func copyPopulation(pop *population) *population {
	newpop := &population{
		perms:  make([]permutation, len(pop.perms)),
		points: copyPointSet(pop.points),
	}
	for i := range newpop.perms {
		newpop.perms[i] = copyPerm(pop.perms[i])
	}

	return newpop
}

func (pop *population) copyPopulation() *population {
	return copyPopulation(pop)
}
