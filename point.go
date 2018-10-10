package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// point is an ordered tuple of values.
type point []float64 // x = (x0, x1, ...)

// pointSet is a set of points.
type pointSet []point

// getPointsFromCSV returns a *pointSet read from a CSV file.
func getPointsFromCSV(filename string) *pointSet {
	var (
		err    error
		pntSet pointSet
		pnt    point
		line   []string
		reader *csv.Reader
		file   *os.File
	)
	file, err = os.Open(filename)
	reader = csv.NewReader(bufio.NewReader(file))
	for {
		line, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		pnt = make(point, len(line))
		for i := range line {
			pnt[i], err = strconv.ParseFloat(line[i], 64)
			if err != nil {
				log.Fatal(err)
			}
		}
		pntSet = append(pntSet, pnt)
	}
	return &pntSet
}
