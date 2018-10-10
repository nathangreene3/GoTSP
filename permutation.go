package main

import "fmt"

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
	q := make(permutation, n)
	k := -1
	for i := n - 2; 0 <= i; i-- {
		if p[i] < p[i+1] {
			k = i
			break
		}
	}
	fmt.Println("K:", k)

	if k == -1 {
		// Generate the first permuation.
		for i := 0; i < len(p); i++ {
			q[i] = i
		}
	} else {
		// Find largest index j > k such that p[k] < p[j].
		// If k = -1, return the base permutation
		// (0, 1, 2, ..., n-1).
		j := -1
		for i := n - 1; k < i; i-- { // Note k cannot ever be n-1
			if p[k] < p[i] {
				j = i
				break
			}
		}
		fmt.Println("J:", j)

		// Copy p[0:k].
		for i := 0; i < k; i++ {
			q[i] = p[i]
		}

		// Swap p[k] and p[j].
		q[k] = p[j]
		q[j] = p[k]

		// 4. Reverse p[k+1:] ignoring p[j].
		for i := k + 1; i < j; i++ {
			q[i] = p[n-i-1]
		}
		for i := j + 1; i < n; i++ {
			q[i] = p[n-i-1]
		}
	}
	return q
}

// basePerm returns the base permutation (012...n-1).
func basePerm(n int) *permutation {
	p := make(permutation, n)
	for i := range p {
		p[i] = i
	}
	return &p
}

// heapPermute is taken from Analysis of Algorithms, page 54.
func heapPermute(a permutation, n int) permutation {
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
