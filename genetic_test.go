package main

import "testing"

func TestCross(t *testing.T) {
	n := 4
	p, q := basePerm(n), basePerm(n)
	for i := 0; i < 10; i++ {
		p, q = cross(p, q)
		if !isPerm(p) || !isPerm(q) {
			t.Fatalf("cross produced non-permutation(s): %v and %v\n", p, q)
		}
	}
}

func TestMutate(t *testing.T) {
	p := basePerm(4)
	for i := 0; i < 10; i++ {
		p = mutate(p)
		if !isPerm(p) {
			t.Fatalf("mutate produced non-permutation: %v\n", p)
		}
	}
}

func TestReproduce(t *testing.T) {

}
func TestCopyPermutation(t *testing.T) {
	p := basePerm(4)
	for i := 0; i < 10; i++ {
		p = copyPerm(p)
		if !isPerm(p) {
			t.Fatalf("copyPermutation produced non-permutation: %v\n", p)
		}
	}
}

func TestCopyPoint(t *testing.T) {

}
func TestCopyPointSet(t *testing.T) {

}
