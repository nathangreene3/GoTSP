package main

// permutation is an ordering on Z(n).
type permutation []int

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

	q := make(permutation, n)
	copy(q, p)

	// Swap p[k] and p[j].
	q[k], q[j] = q[j], q[k]

	// Reverse p[k+1:].
	a := k + 1
	b := n - 1
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

func nextPerm2(p permutation) {
	facts := make([]int, len(p))
	for i := range facts {
		facts[i] = factorial(i)
	}
	p[len(p)-1] = 0
	for i := 0; i < len(p)-1; i++ {
		// d:=
	}
}
