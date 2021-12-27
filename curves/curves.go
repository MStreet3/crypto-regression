package curves

import (
	"log"
	"math"
)

func LogCurve(x []float64) float64 {
	if len(x) != 3 {
		log.Fatalf("must have 3 dimensions got %d", len(x))
	}
	xs := x[0]
	a := x[1]
	b := x[2]
	return a*math.Log(xs) + b
}

func LinearCurve(x []float64) float64 {
	if len(x) != 3 {
		log.Fatalf("must have 3 dimensions got %d", len(x))
	}
	xs := x[0]
	a := x[1]
	b := x[2]
	return a*(xs) + b
}
