package objectives

import (
	"math"

	"gonum.org/v1/plot/plotter"
)

type Func func(x []float64) float64

/* makes an objective function that computes the
sum of squared differences between a curve estimate
and an actual value */
func MakeSumSquaresObj(data plotter.XYs, curve Func) Func {
	obj := func(x []float64) float64 {
		var sum float64
		for _, pi := range data {
			values := append([]float64{pi.X}, x...)
			predicted := curve(values)
			sum += math.Pow(pi.Y-predicted, 2)
		}
		return sum
	}
	return obj
}
