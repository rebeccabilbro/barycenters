package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type argT struct {
	Filename string `cli:"filename" usage:"The file to read from"`
}

// Handle possible errors by aborting
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// Close a file opened elsewhere
func closeFile(fi *os.File) {
	err := fi.Close()
	handle(err)
}

// MassPoint is a body with position and mass information
type MassPoint struct {
	x, y, z, mass float64
}

func addMassPoints(a MassPoint, b MassPoint) MassPoint {
	return MassPoint{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
		a.mass + b.mass,
	}
}

func toWeightedHyperspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.z * a.mass,
		a.mass,
	}
}

func main() {
	// Check arguments
	if len(os.Args) <= 1 {
		fmt.Printf("Not enough arguments! Usage: %s FILENAME\n", os.Args[0])
		os.Exit(1)
	}

	// Open the input file
	file, err := os.Open(os.Args[1])
	handle(err)
	defer closeFile(file)

	// A buffer for the MassPoints
	var masspoints []MassPoint

	// Scan the file for MassPoints
	start_loading := time.Now()
	var count int
	for {
		var x, y, z, mass float64
		_, err = fmt.Fscanf(file, "%f:%f:%f:%f", &x, &y, &z, &mass)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		masspoints = append(masspoints, MassPoint{x, y, z, mass})

		count++
	}

	fmt.Printf("Loaded %d values from file in %s.\n", count, time.Since(start_loading))
	if count <= 1 {
		handle(errors.New("Insufficient number of values; there must be at least one "))
	}

	start_calculation := time.Now()
	// Map the points into the mass-weighed hyperspace
	hyperspace := make([]MassPoint, count)
	for _, v := range masspoints {
		hyperspace = append(hyperspace, toWeightedHyperspace(v))
	}

	// Add up all the points from hyperspace
	var systemHypercenter MassPoint
	for _, v := range hyperspace {
		systemHypercenter = addMassPoints(systemHypercenter, v)
	}

	// Pull the average out of the hyper space into real coordinates
	systemAverage := MassPoint{
		systemHypercenter.x / systemHypercenter.mass,
		systemHypercenter.y / systemHypercenter.mass,
		systemHypercenter.z / systemHypercenter.mass,
		systemHypercenter.mass,
	}

	fmt.Printf("System barycenter is at (%f, %f, %f) and the system's mass is %f.\n",
		systemAverage.x,
		systemAverage.y,
		systemAverage.z,
		systemAverage.mass)
	fmt.Printf("Calculation took %s.\n", time.Since(start_calculation))
}
