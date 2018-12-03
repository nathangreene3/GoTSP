package main

import (
	"math/rand"
	"testing"
)

func TestNilPops(t *testing.T) {
	if !comparePopulations(nil, nil) {
		t.Fatal("explicit nil isn't equal to nil")
	}
	p, q := &population{}, &population{}
	p, q = nil, nil
	if !comparePopulations(p, q) {
		t.Fatal("nil pointers aren't equal")
	}
}

func TestCross(t *testing.T) {
	n := 4
	p, q := basePermutation(n), basePermutation(n)
	for i := 0; i < 10; i++ {
		p, q = cross(p, q)
		if !isPermutation(p) || !isPermutation(q) {
			t.Fatalf("cross produced non-permutation(s): %v and %v\n", p, q)
		}
	}
}

func TestMutate(t *testing.T) {
	// p := basePermutation(4)
	// for i := 0; i < 10; i++ {
	// 	p = mutate(p, nil)
	// 	if !isPermutation(p) {
	// 		t.Fatalf("mutate produced non-permutation: %v\n", p)
	// 	}
	// }
}

func TestReproduce(t *testing.T) {

}
func TestCopyPermutation(t *testing.T) {
	p := basePermutation(4)
	for i := 0; i < 10; i++ {
		p = copyPermutation(p)
		if !isPermutation(p) {
			t.Fatalf("copyPermutation produced non-permutation: %v\n", p)
		}
	}
}

func TestCopyPoint(t *testing.T) {
	n := 10
	p := make(point, n)
	q := make(point, n)
	for i := 0; i < n; i++ {
		for j := range p {
			p[j] = rand.Float64()
		}
		copy(q, p)
		if !comparePoints(p, q) {
			t.Fatalf("points are not equal: %v %v\n", p, q)
		}
	}
}

func TestCopyPointSet(t *testing.T) {

}

func TestIsPermutation(t *testing.T) {
	n := 10
	p := basePermutation(n)
	if !isPermutation(p) {
		t.Fatalf("isPermutation failed to identify base permutation as a valid permutation: %v\n", p)
	}
	for i := 0; i < n; i++ {
		p = randPermutation(n)
		if !isPermutation(p) {
			t.Fatalf("isPermutation failed to identify random permutation as a valid permutation: %v\n", p)
		}
	}
}
