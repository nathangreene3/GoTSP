package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// point is an ordered tuple of values.
type point []float64 // x = (x0, x1, ...)

// pointSet is a set of points.
type pointSet []point

// importPoints returns a point set read from a CSV file.
func importPoints(filename string) (pointSet, error) {
	if !strings.HasSuffix(strings.ToLower(filename), ".csv") {
		filename += ".csv"
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	var ps pointSet
	var p point
	var line []string

	for {
		line, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				return ps, nil
			}
			return nil, err
		}

		p = make(point, len(line))
		for i := range line {
			p[i], err = strconv.ParseFloat(line[i], 64)
			if err != nil {
				return nil, err
			}
		}

		ps = append(ps, p)
	}
}

// comparePoints returns true if the points are equal in length and are equal for
// each indexed value. It returns false if otherwise.
func comparePoints(p, q point) bool {
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

// comparePointSets returns true if the point sets are equal in length and are
// equal for each indexed point. It returns false if otherwise.
func comparePointSets(p, q pointSet) bool {
	if len(p) != len(q) {
		return false
	}

	for i := range p {
		if !comparePoints(p[i], q[i]) {
			return false
		}
	}

	return true
}

// copyPoint returns a new copy of a given point.
func copyPoint(p point) point {
	q := make(point, len(p))
	copy(q, p)
	return q
}

// copyPoint returns a new copy of a given point.
func (p point) copyPoint() point {
	return copyPoint(p)
}

// copyPointSet returns a new copy of a point set.
func copyPointSet(ps pointSet) pointSet {
	newps := make(pointSet, 0, len(ps))
	for i := range ps {
		newps = append(newps, copyPoint(ps[i]))
	}

	return newps
}

// copyPointSet returns a new copy of a point set.
func (ps pointSet) copyPointSet() pointSet {
	return copyPointSet(ps)
}

// sqDist returns the squared distance between two points. Assumes the points are
// equal in dimension.
func (p point) sqDist(q point) float64 {
	var d, t float64
	for i := range p {
		t = q[i] - p[i]
		d += t * t
	}

	return d
}

// sqDist returns the squared distance between two points. Assumes the points are
// equal in dimension.
func sqDist(p0, p1 point) float64 {
	var d, t float64
	for i := range p0 {
		t = p1[i] - p0[i]
		d += t * t
	}

	return d
}

// dist returns the distance between two points.
func (p point) dist(q point) float64 {
	return math.Sqrt(p.sqDist(q))
}

// dist returns the distance between two points.
func dist(p0, p1 point) float64 {
	return math.Sqrt(sqDist(p0, p1))
}

// totalDist returns the total distance traversed across a path of points. A path
// is a permutation of points. Assumes the point set and the permutation are of
// equal dimension.
func totalDist(ps pointSet, perm permutation) float64 {
	n := len(ps)
	d := dist(ps[perm[0]], ps[perm[n-1]])
	for i := 0; i < n-1; i++ {
		d += dist(ps[perm[i]], ps[perm[i+1]])
	}

	return d
}

// totalSqDist returns the total squared distance traversed across a path of
// points. A path is a permutation of points. Assumes the point set and the
// permutation are of equal dimension.
func totalSqDist(ps pointSet, perm permutation) float64 {
	n := len(ps)
	d := sqDist(ps[perm[0]], ps[perm[n-1]])
	for i := 0; i < n-1; i++ {
		d += sqDist(ps[perm[i]], ps[perm[i+1]])
	}

	return d
}

// lineThrough returns a function that is the line passing through two points.
// Assumes the points are two dimensional.
func lineThrough(x0, x1 point) func(x float64) float64 {
	return func(x float64) float64 { return (x1[1]-x0[1])*(x-x0[0])/(x1[0]-x0[0]) + x0[1] }
}

// crossAt returns the x value that is the crossing point for two lines generated
// by the given points x0, ..., x3. Line01 passes through Line23 only if they do
// not share the same slope and are in fact distinct lines.
func crossAt(x0, x1, x2, x3 point) float64 {
	d01 := diff(x0, x1)
	d02 := diff(x0, x2)
	d23 := diff(x2, x3)
	return (d02[1] + x2[0]*(d23[1]/d23[0]) - x0[0]*(d01[1]/d01[0])) / (d23[1]/d23[0] - d01[1]/d01[0])
}

// diff returns a point where each value is the component- wise difference of the
// points x and y.
func diff(x, y point) point {
	n := len(x)
	d := make(point, 0, n)
	for i := 0; i < n; i++ {
		d = append(d, y[i]-x[i])
	}

	return d
}

// pathsCross determines if the paths (x0,x1) and (x2,x3) cross and returns the
// cross point x if they do.
func pathsCross(x0, x1, x2, x3 point) (float64, bool) {
	x := crossAt(x0, x1, x2, x3)
	if isBetween(x, x0[0], x1[0]) && isBetween(x, x2[0], x3[0]) {
		return x, true
	}

	return x, false
}

// isBetween returns true if x is between a and b.
func isBetween(x, a, b float64) bool {
	if a < b {
		return a <= x && x <= b // Potentially, a == b == x
	}

	return b <= x && x <= a
}
