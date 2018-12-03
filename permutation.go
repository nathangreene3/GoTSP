package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// permutation is an ordering on Z(n).
type permutation []int

// population is a set of permutations.
type population struct {
	perms  []permutation // An ordering of points
	points pointSet      // Set of points
}

// nextPermutation returns the next lexicographical permutation. If
// the current permutation is the last permutation, then the
// base permutation is returned.
func nextPermutation(p permutation) permutation {
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
		return basePermutation(n)
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
	q := copyPermutation(p)
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

// basePermutation returns the base permutation (012...n-1).
func basePermutation(n int) permutation {
	p := make(permutation, n)
	for i := range p {
		p[i] = i
	}
	return p
}

// randPermutation returns a random permutation of length n.
func randPermutation(n int) permutation {
	p := basePermutation(n)
	var a, b int
	for i := 0; i < 3*n; i++ {
		a, b = rand.Intn(n), rand.Intn(n)
		p[a], p[b] = p[b], p[a]
	}
	return p
}

// copyPermutation returns a new permution copy.
func copyPermutation(p permutation) permutation {
	q := make(permutation, len(p))
	copy(q, p)
	return q
}

// compareIntSlices returns true if the two slices are nil
// or if they share the same length and values at each index.
func comparePermutations(p, q permutation) bool {
	if len(p) != len(q) {
		return false
	}
	for i := range p {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

// isPermutation determines if an []int is a permutation.
func isPermutation(a []int) bool {
	m := make(map[int]int)
	for i := range a {
		m[a[i]] = 0
	}
	for i := range a {
		m[a[i]]++
	}
	for _, v := range m {
		if v != 1 {
			return false
		}
	}
	return true
}

// index returns the first index of a value in a given permutation. If
// not found, -1 is returned.
func (p permutation) index(v int) int {
	for i := range p {
		if p[i] == v {
			return i
		}
	}
	return -1
}

// cross returns two new permutions each leading with the values of x and
// y but with trailing values of y and x. Partially mapped crossover
// (PMX) is used to ensure the returned permutations are actual
// permutations. The pivot point is selected at random. For example,
// [1,2,3,4] and [5,6,7,8] might be crossed at index 1 returning
// [1,6,7,8] and [5,2,3,4].
// Source: http://user.ceng.metu.edu.tr/~ucoluk/research/publications/tspnew.pdf
func cross(p, q permutation) (permutation, permutation) {
	u := copyPermutation(p)
	v := copyPermutation(q)
	pivot := rand.Intn(len(p))
	for i := 0; i <= pivot; i++ {
		u[u.index(q[i])] = u[i]
		u[i] = q[i]
		v[v.index(p[i])] = v[i]
		v[i] = p[i]
	}
	return u, v
}

// mutate reverses a subsequence within a permutation. For example,
// [1,2,3,4] might be mutated as [1,3,2,4], or even [3,2,1,4].
func mutate(p permutation, ps pointSet) permutation {
	q := copyPermutation(p)
	b := rand.Intn(len(p)-1) + 1 // 0 < b < n
	a := rand.Intn(b)            // 0 <= a < b
	for a < b {
		q[a], q[b] = q[b], q[a]
		a++
		b--
	}
	return q
}

// reproduce
func reproduce(pop *population, topPct, mutationRate float64) *population {
	nextGen := copyPopulation(pop)
	n := len(nextGen.perms)
	for i := 0; i < int(topPct*float64(n)); i += 2 {
		nextGen.perms[n-i-1], nextGen.perms[n-i-2] = cross(nextGen.perms[i], nextGen.perms[i+1])
		if rand.Float64() < mutationRate {
			nextGen.perms[n-i-1] = mutate(nextGen.perms[n-i-1], nextGen.points)
			nextGen.perms[n-i-2] = mutate(nextGen.perms[n-i-2], nextGen.points)
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
		pop.perms[i] = randPermutation(n)
	}
	return pop
}

func (pop *population) Len() int {
	return len(pop.perms)
}

func (pop *population) Less(i, j int) bool {
	return totalSqDist(pop.points, pop.perms[i]) < totalSqDist(pop.points, pop.perms[j])
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
		newpop.perms[i] = copyPermutation(pop.perms[i])
	}
	return newpop
}

func populationToString(pop *population, name string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s addr: %v\n", name, &pop))
	sb.WriteString(fmt.Sprintf("dist: %0.2f\n", totalDist(pop.points, pop.perms[0])))
	sb.WriteString(fmt.Sprintf("perms length: %v\n", len(pop.perms)))
	for i := range pop.perms {
		sb.WriteString(fmt.Sprintf("perm[%d]: %v\n", i, pop.perms[i]))
	}
	return sb.String()
}

func comparePopulations(p, q *population) bool {
	if p == nil {
		if q == nil {
			return true
		}
		return false
	}
	if q == nil {
		return false
	}

	if len(p.perms) != len(q.perms) {
		return false
	}
	for i := range p.perms {
		if !comparePermutations(p.perms[i], q.perms[i]) {
			return false
		}
	}

	if !comparePointSets(p.points, q.points) {
		return false
	}

	return true
}
