package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mstreet3/crypto-regression/curves"
	"github.com/mstreet3/crypto-regression/objectives"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func makePlot(path string, data plotter.XYs, fit plotter.XYs) {

	// xticks defines how we convert and display time.Time values.
	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p := plot.New()
	p.Title.Text = "Weekly Closing Price (LINKUSD)"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "LINK Price ($)"
	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 255, A: 255}

	// p.Add(line, points)
	p.Add(points)

	if fit != nil {
		line, _, err := plotter.NewLinePoints(fit)
		if err != nil {
			log.Panic(err)
		}
		line.Color = color.RGBA{G: 255, A: 255}
		p.Add(line)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Panicf("could not create %s: %v", path, err)
	}
	err = p.Save(10*vg.Centimeter, 5*vg.Centimeter, path)
	if err != nil {
		log.Panic(err)
	}
	if err = f.Close(); err != nil {
		log.Panicf("could not close %s: %v", path, err)
	}
}

func readData(path string) (plotter.XYs, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var xys plotter.XYs
	i := -1
	s := bufio.NewScanner(f)
	for s.Scan() {
		// skip header row
		i++
		if i == 0 {
			continue
		}

		var x, y float64

		// time_period_start;time_period_end;time_open;time_close;price_close
		values := strings.Split(s.Text(), ";")

		time_period_end := values[1]
		price_close := values[4]

		var xtime time.Time
		xtime, err = time.Parse(time.RFC3339Nano, time_period_end)
		if err != nil {
			log.Printf("error parsing time %q: %v", time_period_end, err)
			continue
		}
		x = float64(xtime.Unix())
		y, err = strconv.ParseFloat(price_close, 64)
		if err != nil {
			log.Printf("error parsing closing price %q: %v", price_close, err)
			continue
		}
		xys = append(xys, struct{ X, Y float64 }{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}
	return xys, nil
}

func solveLogarithm() {

	xys, err := readData("7DAY-link-data.csv")
	if err != nil {
		log.Fatalf("error loading data: %v", err.Error())
	}

	predictor := curves.LogCurve
	ss := objectives.MakeSumSquaresObj(xys, predictor)
	leastSquares := optimize.Problem{
		Func: ss,
	}

	initX := []float64{0.1, 0.5}
	result, err := optimize.Minimize(leastSquares, initX, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = result.Status.Err(); err != nil {
		log.Fatal(err)
	}

	var fit plotter.XYs
	for _, pi := range xys {
		pred := predictor([]float64{pi.X, result.X[0], result.X[1]})
		fit = append(fit, struct{ X, Y float64 }{pi.X, pred})
	}
	makePlot("fit-result.png", xys, fit)
}

func main() {
	solveLogarithm()
}
