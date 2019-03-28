package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

// permutation is an ordering on Z(n).
type permutation []int

// population is a set of permutations over a point set.
type population struct {
	perms        []permutation // An ordering of points
	shortestPerm permutation   // Best solution in this population
	points       pointSet      // Set of points
}

// breedFunc returns two new permutations bred from two given permutations.
type breedFunc func(permutation, permutation) (permutation, permutation)

// mutateFunc returns a new permutation altered from the given permutation on a
// given point set.
type mutateFunc func(permutation, pointSet) permutation

func (p permutation) String() string {
	n := len(p)
	sb := strings.Builder{}
	sb.Grow(2 * (n + 1))

	sb.WriteByte('[')
	sb.WriteString(strconv.Itoa(p[0]))
	for i := 1; i < n; i++ {
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(p[i]))
	}
	sb.WriteByte(']')

	return sb.String()
}

// basePermutation returns the base permutation (012...n-1).
func basePermutation(n int) permutation {
	p := make(permutation, 0, n)
	for i := 0; i < n; i++ {
		p = append(p, i)
	}

	return p
}

// nextPermutation returns the next lexicographical permutation. If the current
// permutation is the last permutation, then the base permutation is returned.
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

// randPermutation returns a random permutation of length n.
func randPermutation(n int) permutation {
	p := make(permutation, n)
	var j int
	for i := range p {
		j = rand.Intn(i + 1)
		p[i] = p[j]
		p[j] = i
	}

	return p
}

// copyPermutation returns a new permution copy.
func copyPermutation(p permutation) permutation {
	q := make(permutation, len(p))
	copy(q, p)
	return q
}

// compareIntSlices returns true if the two slices are nil or if they share the
// same length and values at each index.
func equalPermutations(p, q permutation) bool {
	if p == nil {
		return q == nil
	}

	if q == nil {
		return false
	}

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
		m[a[i]]++
		if 1 < m[a[i]] {
			return false
		}
	}

	return true
}

// isBase returns true if a given permutation is the base permutation.
func isBase(p permutation) bool {
	for i := range p {
		if p[i] != i {
			return false
		}
	}

	return true
}

// index returns the first index of a value in a given permutation. If not found,
// -1 is returned.
func (p permutation) index(v int) int {
	for i := range p {
		if p[i] == v {
			return i
		}
	}

	return -1
}

// permutatePointSet returns a new point set ordered by a given permutation. This
// does not sort the point set.
func permutatePointSet(ps pointSet, p permutation) pointSet {
	newps := make(pointSet, len(ps))
	for i := range newps {
		newps[i] = ps[p[i]]
	}

	return newps
}

// pmx returns two new permutions each leading with the values of x and y but with
// trailing values of y and x. Partially mapped crossover (PMX) is used to
// ensure the returned permutations are actual permutations. The pivot point is
// selected at random.
func pmx(p, q permutation) (permutation, permutation) {
	// Source: http://user.ceng.metu.edu.tr/~ucoluk/research/publications/tspnew.pdf

	u := copyPermutation(p)
	v := copyPermutation(q)
	pivot := rand.Intn(len(p))
	var t int

	for i := 0; i <= pivot; i++ {
		t = q[i]
		u[u.index(t)] = u[i]
		u[i] = t

		t = p[i]
		v[v.index(t)] = v[i]
		v[i] = t
	}

	return u, v
}

// reverseSubsequence reverses a subsequence within a permutation. For example, [1,2,3,4]
// might be mutated as [1,3,2,4], or even [3,2,1,4].
func reverseSubsequence(p permutation, ps pointSet) permutation {
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

// reproduce returns a new population as the next generation. The top percent of
// the population gets to reproduce with the possibility of a mutation occurring
// at a given rate. The population size will remain constant. The next generation
// will be returned sorted.
func reproduce(pop *population, topPct, mutationRate float64, f mutateFunc, g breedFunc) *population {
	nextGen := copyPopulation(pop)
	n := len(nextGen.perms)
	topn := int(topPct * float64(n))

	for i := 0; i < topn; i += 2 {
		nextGen.perms[n-i-1], nextGen.perms[n-i-2] = pmx(nextGen.perms[i], nextGen.perms[i+1])
		if rand.Float64() < mutationRate || equalPermutations(nextGen.perms[n-i-1], nextGen.perms[n-i-2]) {
			nextGen.perms[n-i-1] = f(nextGen.perms[n-i-1], nextGen.points)
			nextGen.perms[n-i-2] = f(nextGen.perms[n-i-2], nextGen.points)
		}
	}

	sort.Sort(nextGen)
	if totalDist(nextGen.points, nextGen.perms[0]) < totalDist(nextGen.points, nextGen.shortestPerm) {
		nextGen.shortestPerm = copyPermutation(nextGen.perms[0])
	}

	return nextGen
}

// randPopulation returns a random population of a given size over a set
// of points.
func randPopulation(size int, ps pointSet) *population {
	if ps == nil || len(ps) == 0 {
		panic("pointSet must have at least one point")
	}
	if size < 1 {
		panic("population size must be positive")
	}

	pop := &population{
		perms:  make([]permutation, 0, size),
		points: copyPointSet(ps),
	}

	n := len(ps)
	for i := 0; i < size; i++ {
		pop.perms = append(pop.perms, rand.Perm(n))
	}

	sort.Sort(pop)
	pop.shortestPerm = copyPermutation(pop.perms[0])

	return pop
}

// copyPopulation returns a copy of a given population.
func copyPopulation(pop *population) *population {
	newpop := &population{
		perms:        make([]permutation, 0, len(pop.perms)),
		shortestPerm: copyPermutation(pop.shortestPerm),
		points:       copyPointSet(pop.points),
	}

	for i := 0; i < pop.Len(); i++ {
		newpop.perms = append(newpop.perms, copyPermutation(pop.perms[i]))
	}

	return newpop
}

// populationToString returns a formatted string representation of a given
// population.
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

// equalPopulations determines if two populations contain equal values.
func equalPopulations(p, q *population) bool {
	if p == nil {
		return q == nil
	}

	if q == nil {
		return false // p is not nil, q is
	}

	if len(p.perms) != len(q.perms) {
		return false
	}

	for i := range p.perms {
		if !equalPermutations(p.perms[i], q.perms[i]) {
			return false
		}
	}

	if !equalPermutations(p.shortestPerm, q.shortestPerm) {
		return false
	}

	if !comparePointSets(p.points, q.points) {
		return false
	}

	return true
}

// TODO
func (pop *population) distinctMembers() int {
	n := len(pop.perms)
	if n == 0 {
		return 0
	}

	c := make([]int, n)
	var count int
	for i := 0; i+1 < n; i++ {
		for j := i + 1; j < n; j++ {
			if equalPermutations(pop.perms[i], pop.perms[j]) {
				c[i]++
			}
		}
	}

	for i := range c {
		if c[i] == 0 {
			count++
		}
	}
	return count
}

// importPermutation retrieves a permutation from a given csv file.
func importPermutation(filename string) (permutation, error) {
	filename = strings.ToLower(filename)
	if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
		filename += ".csv"
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
		}

		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	var p permutation
	var line []string
	var v int64

	for {
		line, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				return p, nil
			}
			return nil, err
		}

		p = make(permutation, 0, len(line))
		for i := range line {
			v, err = strconv.ParseInt(line[i], 10, 64)
			if err != nil {
				return nil, err
			}

			p = append(p, int(v))
		}
	}
}

// export writes a permutation to a given file. The file contents will be
// overwritten if it exists and will be created if it doesn't exist.
func (p permutation) exportPermutation(filename string) error {
	filename = strings.ToLower(filename)
	if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
		filename += ".csv"
	}

	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	output := make([]string, 0, len(p))
	for i := range p {
		output = append(output, strconv.Itoa(p[i]))
	}

	writer := csv.NewWriter(file)
	err = writer.Write(output)
	if err != nil {
		return err
	}
	writer.Flush()

	return writer.Error()
}

// Len returns the size of the population (number of permutations).
func (pop *population) Len() int {
	return len(pop.perms)
}

// Less returns true if a permutation i gives a smaller total distance than a
// permutation j.
func (pop *population) Less(i, j int) bool {
	return totalDist(pop.points, pop.perms[i]) < totalDist(pop.points, pop.perms[j])
}

// Swap swaps two permutations i and j.
func (pop *population) Swap(i, j int) {
	pop.perms[i], pop.perms[j] = pop.perms[j], pop.perms[i]
}
