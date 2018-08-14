package main

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "math/rand"
)

func main() {
  // expect two arguments; the binary and the desired # of points
  if len(os.Args) < 2 {
    os.Exit(1)
  }

  // assuming there are enough args, read 2nd as an integer
  nBodies, err := strconv.Atoi(os.Args[1])

  // exit if the user input isn't sensible
  if err != nil {
    os.Exit(1)
  }

  rand.Seed(time.Now().Unix())

  posMax := 100 // maximum deviation from the center of any axis
  massMax := 5 // maximum mass

  // generate the desired points
  for i := 0; i < nBodies; i++ {
    posX :=rand.Intn(posMax * 2) - posMax
    posY :=rand.Intn(posMax * 2) - posMax
    posZ :=rand.Intn(posMax * 2) - posMax
    mass := rand.Intn(massMax - 1) +1
    fmt.Printf("%d:%d:%d:%d\n", posX, posY, posZ, mass)

  }
}
