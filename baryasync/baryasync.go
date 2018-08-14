package main

import (
  "os"
  "io"
	"fmt"
	"time"
  "bufio"
	"errors"
)

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

func avgMassPoints(a MassPoint, b MassPoint) MassPoint {
	sum := addMassPoints(a, b)
	return MassPoint{
		sum.x / 2,
		sum.y / 2,
		sum.z / 2,
		sum.mass,
	}
}

func toWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.z * a.mass,
		a.mass,
	}
}

func fromWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x / a.mass,
		a.y / a.mass,
		a.z / a.mass,
		a.mass,
	}
}

func avgMassPointsWeighted(a MassPoint, b MassPoint) MassPoint {
	aWeighted := toWeightedSubspace(a)
	bWeighted := toWeightedSubspace(b)
	return fromWeightedSubspace(avgMassPoints(aWeighted, bWeighted))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func closeFile(fi *os.File) {
	err := fi.Close()
	handle(err)
}

func stringToPointAsync(s string, c chan<- MassPoint, wg *sync.WaitGroup){
  defer wg.Done()
  var newMassPoint MassPoint
  _, err := fmt.Sscanf(s, "%f:%f:%f:%f", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)
  if err != nil {
    return
  }
  c <- newMassPoint
}

func avgMassPointsWeightedAsync(a MassPoint, b MassPoint, c chan<- MassPoint) {
  c <- avgMassPointsWeighted(a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect number of arguments!")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	handle(err)
	defer closeFile(file)

	var masspoints []MassPoint

	startLoading := time.Now()

  r := bufio.NewReader(file)
  masspointsChan := make(chan MassPoint, 128)
  var wg sync.WaitGroup
	for {
    str, err := r.ReadString("\n")
    if len(str) == 0 || err != nil {
      break
    }

    wg.Add(1)
    go stringToPointAsync(str, masspointsChan, &wg)
	}

  syncChan := make(chan bool)
  go func() {wg.Wait(); syncChan <- false}

  run := true

  for run || len(masspointsChan) > 0 {
    select {
    case value := <- masspointsChan
      masspoints := append(masspoints, masspointsChan)
    case _ = <-syncChan:
      run = false
    }
  }

	fmt.Printf("Loaded %d values from file in %s.\n", len(masspoints), time.Since(startLoading))
	if len(masspoints) <= 1 {
		handle(errors.New("Insufficient values."))
	}

  c := make(chan MassPoint, len(masspoints)/2)

	startCalculation := time.Now()
	for len(masspoints) > 1 {
		var newMasspoints []MassPoint

    goroutines := 0
		for i := 0; i < len(masspoints) - 1; i += 2 {
			go avgMassPointsWeighted(masspoints[i], masspoints[i+1], c)
      goroutines ++
    }

    for i := 0; i < goroutines; i ++ {
      newMasspoints = append(newMasspoints, <-c)
    }

		if len(masspoints) % 2 != 0 {
			newMasspoints = append(newMasspoints, masspoints[len(masspoints)-1])
		}

		masspoints = newMasspoints
	}

	systemAverage := masspoints[0] // this is the last one remaining

	fmt.Printf("System barycenter is at (%f, %f, %f) and the system's mass is %f.\n",
		systemAverage.x,
		systemAverage.y,
		systemAverage.z,
		systemAverage.mass)
  fmt.Printf("Calculation took %s.\n", time.Since(startCalculation))
}
